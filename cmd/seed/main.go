package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/etc"
	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/http/api"
	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/persistence/redis"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(etc.GetLogLevel())
	if err := run(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(_ []string) (err error) {
	log.Info("Starting Go seed")
	config, err := etc.GetConfig()
	if err != nil {
		return
	}

	pool := redis.GetPool(config)
	store := redis.NewStore(pool)
	handler := api.NewHandler(pool, store)
	apiServer := api.NewServer(handler)

	complete := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		captured := <-sigint
		log.WithField("signal", captured.String()).Trace("Trapped os signal")

		apiServer.Shutdown()
		_ = pool.Close()

		close(complete)
	}()

	go func() {
		apiServer.ListenAndServe()
	}()

	<-complete
	return
}
