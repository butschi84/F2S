package routes

import (
	"encoding/json"
	"net/http"
)

// *********************************************************
// get current data from dispatcher
// *********************************************************
func getCurrentDispatcherData(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get current dispatcher data")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(f2shub.F2SDispatcherHub.Pretty().DispatcherFunctions)

	// output := ""
	// output += "Dispatcher Data\n"
	// output += "===============\n"
	// // iterate functions
	// for _, function := range f2shub.F2SConfiguration.Functions.Items {
	// 	output += fmt.Sprintf("%s\n", function.Name)

	// 	// get function target
	// 	target, err := f2shub.F2STargets.GetFunctionTargetByFunctionName(function.Name)
	// 	logging.Error(err)

	// 	output += fmt.Sprintf("Endpoints: %d\n", len(target.ServingPods))
	// 	for _, endpoint := range target.ServingPods {
	// 		output += fmt.Sprintf("=> %s (inflight requests: %d)\n", string(endpoint.Address.IP), len(endpoint.InflightRequests))
	// 	}
	// }

	// w.Write([]byte(output))
}
