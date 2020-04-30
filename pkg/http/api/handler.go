package api

import (
	"encoding/json"
	"net/http"

	"github.com/danielpacak/redis-sentinel-client-go-seed/pkg/persistence"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type handler struct {
	store persistence.Store
}

func NewHandler(store persistence.Store) http.Handler {
	handler := &handler{
		store: store,
	}

	router := mux.NewRouter()
	router.Use(handler.logRequest)

	redisRouter := router.PathPrefix("/redis").Subrouter()
	redisRouter.Methods(http.MethodPost).Path("/key").HandlerFunc(handler.setKey)
	redisRouter.Methods(http.MethodGet).Path("/keys").HandlerFunc(handler.getKeys)

	return router
}

func (h *handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Tracef("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (h *handler) setKey(w http.ResponseWriter, r *http.Request) {
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

func (h *handler) getKeys(w http.ResponseWriter, _ *http.Request) {
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
