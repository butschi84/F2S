package routes

import (
	v1alpha1types "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/queue"
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
	for _, f := range F2SHub.F2SConfiguration.Functions.Items {
		if string(f.Spec.Endpoint) == string(target) {
			logging.Info(fmt.Sprintf("found function %s (%s) for target: %s", f.Name, f.UID, target))
			url := fmt.Sprintf("invoke url: http://%s.f2s-containers:%v%s", f.Name, f.Target.Port, f.Target.Endpoint)
			logging.Info(url)

			return &f, nil
		}
	}
	return nil, errors.New("no matching function found for endpoint")
}

// delete a F2SFunction
func invokeFunction(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to invoke a function")

	// parse uid
	logging.Info("parsing target path from request")
	vars := mux.Vars(r)
	key := "/" + vars["target"]

	// make request obj
	request := queue.F2SRequest{
		UID:    F2SHub.F2SEventManager.GenerateUUID(),
		Path:   "/" + vars["target"],
		Method: "GET",
	}

	// put it into queue
	logging.Info("add request to queue")
	F2SHub.F2SQueue.AddRequest(request)

	// wait for completion

	json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("success: %s", key)})

	// also grab the query parameters
	// queryParameters := r.URL.Query()
	// queryString := queryParameters.Encode()

	// // find relevant function for this target
	// f, err := findF2SFunctionForTarget(key)
	// if err != nil {
	// 	logging.Info(fmt.Sprintf("function not found for endpoint %s", key))
	// 	json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("failed - function not found for endpoint %s", key)})
	// 	return
	// }

	// // send 'function invoked' event
	// F2SHub.F2SEventManager.Publish(eventmanager.Event{
	// 	UID:      F2SHub.F2SEventManager.GenerateUUID(),
	// 	Data:     key,
	// 	Function: *f,
	// 	Type:     eventmanager.Event_FunctionInvoked,
	// })

	// // start time measurement
	// start := time.Now()

	// // invoke
	// url := fmt.Sprintf("http://%s.f2s-containers:%v%s?%s", f.Name, f.Target.Port, f.Target.Endpoint, queryString)
	// result, err := httpGet(url)

	// // measure time elapsed
	// elapsed := time.Since(start)
	// logging.Info("Function execution time: %s\n", fmt.Sprintf("%s", elapsed))

	// // send invocation end event
	// F2SHub.F2SEventManager.Publish(eventmanager.Event{
	// 	UID:      F2SHub.F2SEventManager.GenerateUUID(),
	// 	Data:     elapsed,
	// 	Function: *f,
	// 	Type:     eventmanager.Event_FunctionInvokationEnded,
	// })

	// // send results
	// if err != nil {
	// 	logging.Error(fmt.Errorf("error during invocation"))
	// 	logging.Error(err)
	// 	json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("error during invocation: %s", err)})
	// } else {
	// 	logging.Info(fmt.Sprintf("invocation of function %s completed in %v ms", *&f.Name, elapsed.Milliseconds()))
	// 	json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("success: %s", result)})
	// }
}
