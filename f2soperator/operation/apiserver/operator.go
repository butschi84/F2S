package apiserver

import (
	"encoding/json"
	"net/http"
)

func getOperatorState(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get current operator state")

	// set response headers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(f2shub.F2SOperatorState)
}
