package routes

import (
	v1alpha1types "butschi84/f2s/configuration/api/types/v1alpha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func httpGet(url string) (string, error) {
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer response.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Print response body
	return string(body), nil
}

// search k8s crd's (f2sfunction) for matching target
func findF2SFunctionForTarget(target string) (*v1alpha1types.Function, error) {
	for _, f := range F2SConfiguration.Functions.Items {
		if string(f.Spec.Endpoint) == string(target) {
			logging.Println(fmt.Sprintf("found function %s (%s) for target: %s", f.Name, f.UID, target))
			url := fmt.Sprintf("invoke url: http://%s.f2s-containers:%v%s", f.Name, f.Target.Port, f.Target.Endpoint)
			logging.Println(url)

			return &f, nil
		}
	}
	return nil, errors.New("no matching function found for endpoint")
}

// delete a F2SFunction
func invokeFunction(w http.ResponseWriter, r *http.Request) {
	logging.Println("request to invoke a function")

	// parse uid
	logging.Println("parsing target path from request")
	vars := mux.Vars(r)
	key := "/" + vars["target"]

	f, err := findF2SFunctionForTarget(key)
	if err != nil {
		logging.Println(fmt.Sprintf("function not found for endpoint %s", key))
		json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("failed - function not found for endpoint %s", key)})
		return
	}

	// invoke
	url := fmt.Sprintf("http://%s.f2s-containers:%v%s", f.Name, f.Target.Port, f.Target.Endpoint)
	result, err := httpGet(url)
	if err != nil {
		logging.Println("error during invocation", err)
		json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("error during invocation: %s", err)})
		return
	}

	json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("success: %s", result)})
}
