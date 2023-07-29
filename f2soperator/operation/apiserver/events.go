package apiserver

import (
	"encoding/json"
	"net/http"
)

// *********************************************************
// get latest events
// *********************************************************
func getLatestEvents(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get all k8s services")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(f2shub.F2SEventManager.LastEvents)
}
