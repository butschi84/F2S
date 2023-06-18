package routes

import (
	"butschi84/f2s/configuration"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// pointer to F2SConfiguration
var F2SConfiguration configuration.F2SConfiguration

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome homepage")
	fmt.Println("endpoint hit: homepage")
}
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	log.Println("request to get all functions")
	json.NewEncoder(w).Encode(F2SConfiguration.Functions)

	// test
	configuration.GetCRDs()
}
func getFunction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("searching for key ", key)
	for _, function := range F2SConfiguration.Functions {
		if function.Metadata.Id == key {
			json.NewEncoder(w).Encode(function)
			return
		}
	}
	fmt.Fprintf(w, "{}")
}

func HandleRequests(config configuration.F2SConfiguration) {
	F2SConfiguration = config
	router := mux.NewRouter().StrictSlash(true)

	// retrieve configured f2s functions
	router.HandleFunc("/functions", returnAllFunctions)
	router.HandleFunc("/functions/{id}", getFunction)
	router.HandleFunc("/", homepage)
	http.ListenAndServe("localhost:8000", router)
}
