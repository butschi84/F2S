package routes

import (
	v1alpha1types "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/queue"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
	for _, f := range f2shub.F2SConfiguration.Functions.Items {
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
		UID:           f2shub.F2SEventManager.GenerateUUID(),
		Path:          "/" + vars["target"],
		Method:        "GET",
		ResultChannel: make(chan queue.F2SRequestResult),
	}

	// put it into queue
	logging.Info("add request to queue")
	f2shub.F2SQueue.AddRequest(request)

	// wait for completion
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	select {
	case result := <-request.ResultChannel:
		logging.Info(fmt.Sprintf("Request completed: %s", result.Result))
		json.NewEncoder(w).Encode(result)
	case <-ctx.Done():
		fmt.Println("Request Timeout reached, cancelling goroutine")
		json.NewEncoder(w).Encode(Status{Status: fmt.Sprintf("failed: %s", key)})
	}
}
