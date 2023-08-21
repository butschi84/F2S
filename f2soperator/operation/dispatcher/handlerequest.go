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
			Result:  map[string]interface{}{},
			Details: fmt.Sprintf("function %s: 'request_timeout' after %dms", req.Path, timeout),
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
	logging.Info(fmt.Sprintf("[%s] processing invocation request: %s", req.UID, req.Path))

	// find function target
	logging.Debug(fmt.Sprintf("[%s] search function target for endpoint: %s", req.UID, req.Path))
	functionTarget, err := f2shub.F2SDispatcherHub.GetDispatcherFunctionByEndpoint(req.Path)
	if err != nil {
		logging.Error(fmt.Errorf("[%s] cannot serve request. function target not found for endpoint %s", req.UID, req.Path))
		logging.Error(err)
	}
	logging.Debug(fmt.Sprintf("[%s] function target is: %s", req.UID, functionTarget.Function.Name))

	// send 'function invoked' event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:         f2shub.F2SEventManager.GenerateUUID(),
		Data:        functionTarget.Function.Prettify(),
		Type:        eventmanager.Event_FunctionInvoked,
		Description: fmt.Sprintf("[%s] function %s has been invoked", req.UID, functionTarget.Function.Name),
	})

	// prepare output
	requestResult := queue.F2SRequestResult{
		Result:                     map[string]interface{}{},
		Success:                    false,
		UID:                        req.UID,
		Duration:                   0.0,
		DurationPerInflightRequest: 0.0,
		Request:                    req,
	}

	// wait for function pod to be available
	// => operator will scale up from 0 to 1 after receiving event 'Event_FunctionInvoked'
	if len(functionTarget.ServingPods) == 0 {
		err := waitForTargetPod(functionTarget)
		if err != nil {
			logging.Error(fmt.Errorf("[%s] aborting function '%s'. scale from 0 failed: %s", req.UID, functionTarget.Function.Name, err.Error()))
			// send result to channel
			requestResult.Details = fmt.Sprintf("[%s] aborting function '%s'. scale from 0 failed: %s", req.UID, functionTarget.Function.Name, err.Error())
			requestResult.Success = false
			*result <- requestResult
			return
		}
	}

	// get the pod that will actually serve the request
	pod, err := functionTarget.ServeRequest(req)
	if err != nil {
		logging.Error(fmt.Errorf("[%s] cannot serve request because cannot determine which pod should serve the request: %s", req.UID, err.Error()))
		requestResult.Details = fmt.Sprintf("[%s] aborting function '%s' invocation because target cannot serve request: %s", req.UID, functionTarget.Function.Name, err.Error())
		requestResult.Success = false
		*result <- requestResult
		return
	}

	// maybe check here if pod is ready ?

	// start time measurement
	start := time.Now()

	// invoke function on target pod
	var requestErr error
	url := fmt.Sprintf("http://%s:%v%s", string(pod.Address.IP), functionTarget.Function.Target.Port, functionTarget.Function.Target.Endpoint)
	// url := fmt.Sprintf("http://127.0.0.1:59514")
	logging.Info(fmt.Sprintf("[%s] request url: %s", req.UID, url))
	switch req.Method {
	case "GET":
		requestErr = httpGet(url, &requestResult)
	case "POST":
		requestErr = httpPost(url, req.Payload, &requestResult)
	case "PUT":
		requestErr = httpPut(url, req.Payload, &requestResult)
	case "DELETE":
		requestErr = httpDelete(url, &requestResult)
	default:
		requestErr = httpGet(url, &requestResult)
	}

	if requestErr != nil {
		logging.Error(err)
		// send result to channel
		requestResult.Details = fmt.Sprintf("[%s] error on function http invocation: %s", req.UID, requestErr.Error())
		requestResult.Success = false
	} else {
		// measure time elapsed
		elapsed := time.Since(start).Milliseconds()
		elapsedPerInflight := int(time.Since(start).Milliseconds()) / len(pod.InflightRequests)
		logging.Info("[%s] Function execution time: %s\n", req.UID, fmt.Sprintf("%vms", elapsed))
		logging.Info("[%s] Function execution time per inflight request: %sms\n", req.UID, fmt.Sprintf("%v", elapsedPerInflight))

		// prepare output
		requestResult.Success = true
		requestResult.Duration = float64(elapsed)
		requestResult.DurationPerInflightRequest = float64(elapsedPerInflight)

	}

	// send result to channel
	*result <- requestResult
}
