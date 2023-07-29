package apiserver

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
	config "butschi84/f2s/state/configuration"
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func root(w http.ResponseWriter, r *http.Request) {
	logging.Info("endpoint hit: homepage")
	json.NewEncoder(w).Encode(Status{Status: "up"})
}

// *********************************************************
// all functions
// *********************************************************
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get all functions")

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

	functions := config.ActiveConfiguration.Functions

	logging.Info("searching for uid: ", key)
	for _, function := range functions.Items {
		if string(function.ObjectMeta.UID) == key {
			json.NewEncoder(w).Encode(function.Prettify())
			return
		}
	}

	fmt.Fprintf(w, "{}")
}

// *********************************************************
// create a function
// *********************************************************
func createFunction(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to create a new function")

	logging.Info("parsing request body")
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Define a variable to hold the JSON data
	var function typesV1alpha1.PrettyFunction

	// Unmarshal the JSON data into the function struct
	if err := json.Unmarshal(body, &function); err != nil {
		logging.Error(err)
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}
	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}

	// create f2sfunction crd in k8s
	result, err := kubernetesservice.CreateF2SFunction(&function)

	json.NewEncoder(w).Encode(result)
}

// delete a F2SFunction
func deleteFunction(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to delete a function")

	// parse uid
	logging.Info("parsing uid from request")
	vars := mux.Vars(r)
	key := vars["id"]

	err := kubernetesservice.DeleteF2SFunction(key)
	if err != nil {
		json.NewEncoder(w).Encode(Status{Status: "failed"})
		return
	}

	json.NewEncoder(w).Encode(Status{Status: "success"})
}
