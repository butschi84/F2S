package apiserver

import (
	"encoding/json"
	"net/http"
)

// *********************************************************
// get current configucluster stateration
// *********************************************************
func getClusterState(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get current cluster state")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	jsonBytes, err := json.MarshalIndent(f2shub.F2SClusterState, "", "  ")
	if err != nil {
		return
	}

	_, err = w.Write(jsonBytes)
}
