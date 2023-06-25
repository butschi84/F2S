package routes

import (
	"butschi84/f2s/services/prometheus"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// *********************************************************
// all deployments
// *********************************************************
func getPrometheusMetric(w http.ResponseWriter, r *http.Request) {
	logging.Println("request to get a prometheus metric")

	vars := mux.Vars(r)
	functionName := vars["functionname"]
	metricName := vars["metricname"]

	result, err := prometheus.ReadPrometheusMetric(metricName, map[string]string{"functionname": functionName})
	if err != nil {
		logging.Println(fmt.Sprintf("could not read prometheus metric %s", err))
		json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("could not read prometheus metric %s", err)})
		return
	}

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}
