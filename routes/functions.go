package routes

import (
	config "butschi84/f2s/configuration"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func root(w http.ResponseWriter, r *http.Request) {
	logging.Println("endpoint hit: homepage")
	json.NewEncoder(w).Encode(Status{Status: "up"})
}

// *********************************************************
// all functions
// *********************************************************
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	logging.Println("request to get all functions")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	functions := config.ActiveConfiguration.Functions

	json.NewEncoder(w).Encode(functions.Prettify())
}

// *********************************************************
// specific function
// *********************************************************
func getFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	functions := config.GetCRDs()
	logging.Println("searching for uid: ", key)
	for _, function := range functions.Items {
		if string(function.ObjectMeta.UID) == key {
			json.NewEncoder(w).Encode(function.Prettify())
			return
		}
	}

	fmt.Fprintf(w, "{}")
}
