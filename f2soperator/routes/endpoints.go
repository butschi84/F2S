package routes

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"encoding/json"
	"net/http"
)

// *********************************************************
// all endpoints
// *********************************************************
func getAllEndpoints(w http.ResponseWriter, r *http.Request) {
	logging.Println("request to get all endpoints")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	endpoints, _ := kubernetesservice.ListEndpoints()

	json.NewEncoder(w).Encode(endpoints)
}
