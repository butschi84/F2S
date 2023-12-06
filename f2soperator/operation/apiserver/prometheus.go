package apiserver

import (
	"butschi84/f2s/services/prometheus"
	"encoding/json"
	"fmt"
	"net/http"
)

// *********************************************************
// all deployments
// *********************************************************
func getPrometheusMetric(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get a prometheus metric")

	// Get the "query" parameter from the URL query string
	query := r.URL.Query().Get("query")

	result, err := prometheus.ReadPrometheusMetric(f2shub.F2SConfiguration, query)
	if err != nil {
		logging.Error(fmt.Sprintf("could not read prometheus metric: %s", err))
		json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("could not read prometheus metric %s", err)})
		return
	}

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}
