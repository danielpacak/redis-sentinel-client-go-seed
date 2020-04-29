package api

import (
	"encoding/json"
	"net/http"

	"github.com/danielpacak/redis-ha-seed/pkg/persistence"

	log "github.com/sirupsen/logrus"
)

type handler struct {
	store persistence.Store
}

func NewHandler(store persistence.Store) http.Handler {
	handler := &handler{
		store: store,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.getKeys)
	return mux

}

func (h *handler) getKeys(w http.ResponseWriter, r *http.Request) {
	log.Debug("Handling GET /keys request")
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
