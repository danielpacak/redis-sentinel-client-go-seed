package api

import (
	"encoding/json"
	"net/http"

	xredis "github.com/gomodule/redigo/redis"

	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/persistence"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type handler struct {
	pool  *xredis.Pool
	store persistence.Store
}

func NewHandler(pool *xredis.Pool, store persistence.Store) http.Handler {
	handler := &handler{
		pool:  pool,
		store: store,
	}

	router := mux.NewRouter()
	router.Use(handler.logRequest)

	redisRouter := router.PathPrefix("/redis").Subrouter()
	redisRouter.Methods(http.MethodPost).Path("/key").HandlerFunc(handler.set)
	redisRouter.Methods(http.MethodGet).Path("/keys").HandlerFunc(handler.keys)
	redisRouter.Methods(http.MethodGet).Path("/info").HandlerFunc(handler.info)
	redisRouter.Methods(http.MethodGet).Path("/pool/stats").HandlerFunc(handler.poolStats)

	return router
}

func (h *handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Tracef("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (h *handler) set(w http.ResponseWriter, r *http.Request) {
	jsonRequest := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&jsonRequest)
	if err != nil {
		log.WithError(err).Error("Error while decoding JSON request")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = h.store.Set(jsonRequest.Key, jsonRequest.Value)
	if err != nil {
		log.WithError(err).Error("Error while setting key")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) keys(w http.ResponseWriter, _ *http.Request) {
	keys, err := h.store.Keys()
	if err != nil {
		log.WithError(err).Error("Error while getting keys")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(keys)
	if err != nil {
		log.WithError(err).Error("Error while encoding keys to JSON")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) info(w http.ResponseWriter, _ *http.Request) {
	infos, err := h.store.Info()
	if err != nil {
		log.WithError(err).Error("Error while getting info")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&infos)
	if err != nil {
		log.WithError(err).Error("Error while encoding info to JSON")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) poolStats(w http.ResponseWriter, _ *http.Request) {
	stats := h.pool.Stats()
	err := json.NewEncoder(w).Encode(&stats)
	if err != nil {
		log.WithError(err).Error("Error while encoding stats to JSON")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
