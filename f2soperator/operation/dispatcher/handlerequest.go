package dispatcher

import (
	"butschi84/f2s/state/eventmanager"
	"butschi84/f2s/state/queue"
	"context"
	"fmt"
	"time"
)

// handle function invocations with request_timeout
func handleRequestsWithTimeout(req queue.F2SRequest) {
	timeout := f2shub.F2SConfiguration.Config.F2S.Timeouts.RequestTimeout

	// Create a context with a timeout of 'request_timeout'
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	requestResultChan := make(chan queue.F2SRequestResult)

	// Start the goroutine
	go handleRequest(req, &requestResultChan)

	select {
	case result := <-requestResultChan:
		req.ResultChannel <- result

		// send request completed event
		f2shub.F2SEventManager.Publish(eventmanager.Event{
			UID:         req.UID,
			Data:        result,
			Type:        eventmanager.Event_FunctionInvokationEnded,
			Description: fmt.Sprintf("%s => function %s completed with result: %v", req.UID, result.Request.Path, result.Success),
		})
		return
	case <-ctx.Done():
		logging.Warn(fmt.Sprintf("request for function call '%s' timed out after %dms!", req.Path, timeout))
		// send result to channel
		result := queue.F2SRequestResult{
			Result:  fmt.Sprintf("function %s: 'request_timeout' after %dms", req.Path, timeout),
			Success: false,
			UID:     req.UID,
			Request: req,
		}
		req.ResultChannel <- result

		// send request completed event
		f2shub.F2SEventManager.Publish(eventmanager.Event{
			UID:         req.UID,
			Data:        result,
			Type:        eventmanager.Event_FunctionInvokationEnded,
			Description: fmt.Sprintf("%s => function call ended with result: %v", req.UID, result.Success),
		})
	}
}

// handle function invocations
func handleRequest(req queue.F2SRequest, result *chan queue.F2SRequestResult) {
	logging.Info(fmt.Sprintf("processing invocation request: %s (%s)", req.UID, req.Path))

	// find function target
	logging.Debug(fmt.Sprintf("search function target for endpoint: %s", req.Path))
	functionTarget, err := f2shub.F2SDispatcherHub.GetFunctionTargetByEndpoint(req.Path)
	if err != nil {
		logging.Error(fmt.Errorf("cannot serve request %s. function target not found for endpoint %s", req.UID, req.Path))
		logging.Error(err)
	}
	logging.Debug(fmt.Sprintf("function target is: %s", functionTarget.Function.Name))

	// send 'function invoked' event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:         f2shub.F2SEventManager.GenerateUUID(),
		Data:        functionTarget.Function.Prettify(),
		Type:        eventmanager.Event_FunctionInvoked,
		Description: fmt.Sprintf("function %s has been invoked", functionTarget.Function.Name),
	})

	// wait for function pod to be available
	// => operator will scale up from 0 to 1 after receiving event 'Event_FunctionInvoked'
	if len(functionTarget.ServingPods) == 0 {
		err := waitForTargetPod(functionTarget)
		if err != nil {
			logging.Error(fmt.Errorf("aborting function '%s'. scale from 0 failed: %s", functionTarget.Function.Name, err.Error()))
			// send result to channel
			*result <- queue.F2SRequestResult{
				Result:  fmt.Sprintf("aborting function '%s'. scale from 0 failed: %s", functionTarget.Function.Name, err.Error()),
				Success: false,
				UID:     req.UID,
			}
			return
		}
	}

	// get the pod that will actually serve the request
	pod, err := functionTarget.ServeRequest(req)
	if err != nil {
		logging.Error(err)
		*result <- queue.F2SRequestResult{
			Result:  fmt.Sprintf("aborting function '%s' invocation because target cannot serve request: %s", functionTarget.Function.Name, err.Error()),
			Success: false,
			UID:     req.UID,
		}
		return
	}

	// maybe check here if pod is ready ?

	// start time measurement
	start := time.Now()

	// invoke function on target pod
	var httpResult string
	var requestErr error
	url := fmt.Sprintf("http://%s:%v%s", string(pod.Address.IP), functionTarget.Function.Target.Port, functionTarget.Function.Target.Endpoint)
	switch req.Method {
	case "GET":
		httpResult, requestErr = httpGet(url)
	case "POST":
		httpResult, requestErr = httpPost(url, req.Payload)
	case "PUT":
		httpResult, requestErr = httpPut(url, req.Payload)
	case "DELETE":
		httpResult, requestErr = httpDelete(url)
	default:
		httpResult, requestErr = httpGet(url)
	}

	if requestErr != nil {
		logging.Error(err)
		// send result to channel
		*result <- queue.F2SRequestResult{
			Result:  fmt.Sprintf("error on function http invocation: %s", requestErr.Error()),
			Success: false,
			UID:     req.UID,
			Request: req,
		}
	} else {
		// measure time elapsed
		elapsed := time.Since(start).Milliseconds()
		elapsedPerInflight := int(time.Since(start).Milliseconds()) / len(pod.InflightRequests)
		logging.Info("Function execution time: %s\n", fmt.Sprintf("%vms", elapsed))
		logging.Info("Function execution time per inflight request: %sms\n", fmt.Sprintf("%v", elapsedPerInflight))

		// send result to channel
		*result <- queue.F2SRequestResult{
			Result:                     httpResult,
			Success:                    true,
			UID:                        req.UID,
			Duration:                   float64(elapsed),
			DurationPerInflightRequest: float64(elapsedPerInflight),
			Request:                    req,
		}
	}
}
