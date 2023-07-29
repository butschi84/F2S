package apiserver

import (
	"encoding/json"
	"net/http"
)

// *********************************************************
// get current configuration
// *********************************************************
func getConfiguration(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get current configuration")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	jsonBytes, err := json.MarshalIndent(f2shub.F2SConfiguration, "", "  ")
	if err != nil {
		return
	}

	_, err = w.Write(jsonBytes)
}
