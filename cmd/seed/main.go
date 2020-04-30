package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/etc"
	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/persistence/redis"
	xredis "github.com/gomodule/redigo/redis"

	"github.com/FZambia/sentinel"
	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/http/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	if err := run(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(_ []string) (err error) {
	config, err := etc.GetConfig()
	if err != nil {
		return
	}

	sntnl := &sentinel.Sentinel{
		Addrs:      []string{config.RedisAddr},
		MasterName: config.RedisMaster,
		Dial: func(addr string) (xredis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := xredis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	pool := &xredis.Pool{
		MaxIdle:   5,
		MaxActive: 5,
		Dial: func() (xredis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			log.Debugf("Connecting to Redis master: %v", masterAddr)
			c, err := xredis.Dial("tcp", masterAddr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c xredis.Conn, t time.Time) error {
			log.Debugf("Testing connection on borrow: %v", t)
			if !sentinel.TestRole(c, "master") {
				return errors.New("role check failed")
			} else {
				return nil
			}
		},
	}

	store := redis.NewStore(pool)
	handler := api.NewHandler(store)
	apiServer := api.NewServer(handler)

	complete := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		captured := <-sigint
		log.WithField("signal", captured.String()).Trace("Trapped os signal")

		apiServer.Shutdown()

		close(complete)
	}()

	go func() {
		apiServer.ListenAndServe()
	}()

	<-complete
	return
}
