package routes

import (
	"butschi84/f2s/configuration"
	config "butschi84/f2s/configuration"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// pointer to F2SConfiguration
var F2SConfiguration config.F2SConfiguration

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome homepage")
	logging.Println("endpoint hit: homepage")
}
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	logging.Println("request to get all functions")

	functions := config.ActiveConfiguration.Functions

	json.NewEncoder(w).Encode(functions)
}
func getFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	functions := config.GetCRDs()

	logging.Println("searching for key ", key)
	for _, function := range functions.Items {
		if string(function.ObjectMeta.UID) == key {
			json.NewEncoder(w).Encode(function)
			return
		}
	}
	fmt.Fprintf(w, "{}")
}

func HandleRequests(config configuration.F2SConfiguration) {
	F2SConfiguration = config
	router := mux.NewRouter().StrictSlash(true)

	// openAPI spec
	openAPIHandler := http.FileServer(http.Dir("./static/openapi"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", openAPIHandler))

	// retrieve configured f2s functions
	router.HandleFunc("/functions", returnAllFunctions)
	router.HandleFunc("/functions/{id}", getFunction)
	router.HandleFunc("/", homepage)
	http.ListenAndServe("localhost:8000", router)
}
