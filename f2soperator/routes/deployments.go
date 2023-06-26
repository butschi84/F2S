package routes

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"encoding/json"
	"net/http"
)

// *********************************************************
// all deployments
// *********************************************************
func getAllDeployments(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get all functions")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	deployments, _ := kubernetesservice.GetDeployments()

	json.NewEncoder(w).Encode(deployments)
}

// *********************************************************
// create test deployments
// *********************************************************
func createDeployment(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get all functions")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	deployment, _ := kubernetesservice.CreateDeployment("test", "nginx", map[string]string{}, 0)

	json.NewEncoder(w).Encode(deployment)
}
