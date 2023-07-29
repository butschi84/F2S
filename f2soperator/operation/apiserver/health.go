package apiserver

import (
	"encoding/json"
	"net/http"
)

type SimpleHealthResponse struct {
	Status string `json:"status"`
}

// *********************************************************
// simple health endpoint
// *********************************************************
func checkHealth(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get current health")

	response := SimpleHealthResponse{Status: "ok"}

	json.NewEncoder(w).Encode(response)
}
