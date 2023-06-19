package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// delete a F2SFunction
func invokeFunction(w http.ResponseWriter, r *http.Request) {
	logging.Println("request to invoke a function")

	// parse uid
	logging.Println("parsing target path from request")
	vars := mux.Vars(r)
	key := "/" + vars["target"]

	for _, f := range F2SConfiguration.Functions.Items {
		if f.Spec.Endpoint == key {
			logging.Println(fmt.Sprintf("found function %s (%s) for target: %s", f.Name, f.UID, key))
			json.NewEncoder(w).Encode(f)
		}
	}
	logging.Println(fmt.Sprintf("function not found for endpoint %s", key))

	json.NewEncoder(w).Encode(Status{Status: key})
}
