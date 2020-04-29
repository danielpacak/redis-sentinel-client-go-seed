package api

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
}

func NewServer(handler http.Handler) (server *Server) {
	server = &Server{
		server: &http.Server{
			Handler: handler,
			Addr:    ":8080",
		},
	}
	return
}

func (s *Server) ListenAndServe() {
	log.Infof("Starting API server on %s", ":8080")
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error: %v", err)
	}
	log.Trace("API server stopped listening for incoming connections")
}

func (s *Server) Shutdown() {
	log.Trace("API server shutdown started")
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.WithError(err).Error("Error while shutting down API server")
		return
	}
	log.Trace("API server shutdown completed")
}
