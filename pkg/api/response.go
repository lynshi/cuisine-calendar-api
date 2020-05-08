package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type errorResponse struct {
	Message string `json:"message"`
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, statusCode int, err error) {
	response, _ := json.Marshal(errorResponse{
		Message: err.Error(),
	})

	log.Error().Stack().Err(err).Msg("responding with error")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
