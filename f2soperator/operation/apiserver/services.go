package apiserver

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"encoding/json"
	"net/http"
)

// *********************************************************
// all services
// *********************************************************
func getAllServices(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get all k8s services")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	services, _ := kubernetesservice.ListServices()

	json.NewEncoder(w).Encode(services)
}
