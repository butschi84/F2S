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

type Status struct {
	Status string `json:"status"`
}

func root(w http.ResponseWriter, r *http.Request) {
	logging.Println("endpoint hit: homepage")
	json.NewEncoder(w).Encode(Status{Status: "up"})
}
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	logging.Println("request to get all functions")

	// set response headers
	w.Header().Set("Content-Type", "application/json")

	functions := config.ActiveConfiguration.Functions

	json.NewEncoder(w).Encode(functions.Prettify())
}
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

func HandleRequests(config configuration.F2SConfiguration) {
	F2SConfiguration = config
	router := mux.NewRouter().StrictSlash(true)

	// openAPI spec
	openAPIHandler := http.FileServer(http.Dir("./static/openapi"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", openAPIHandler))

	// retrieve configured f2s functions
	router.HandleFunc("/functions", returnAllFunctions)
	router.HandleFunc("/functions/{id}", getFunction)

	router.HandleFunc("/", root)

	http.ListenAndServe("localhost:8000", router)
}
