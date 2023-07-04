package dispatcher

import (
	"butschi84/f2s/state/eventmanager"
	"butschi84/f2s/state/queue"
	"fmt"
	"time"
)

// handle function invocations
func handleRequests(req queue.F2SRequest) {
	logging.Info("processing new function invocation request")

	// find function target
	functionTarget, err := f2shub.F2STargets.GetFunctionTargetByEndpoint(req.Path)
	if err != nil {
		logging.Error(fmt.Errorf("cannot serve request %s. function target not found for endpoint %s", req.UID, req.Path))
		logging.Error(err)
	}

	// send 'function invoked' event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:  f2shub.F2SEventManager.GenerateUUID(),
		Data: functionTarget.Function.Prettify(),
		Type: eventmanager.Event_FunctionInvoked,
	})

	// wait for function pod to be available
	// => operator will scale up from 0 to 1 after receiving event 'Event_FunctionInvoked'
	if len(functionTarget.ServingPods) == 0 {
		err := waitForTargetPod(functionTarget)
		if err != nil {
			logging.Error(fmt.Errorf("aborting function '%s' invocation because scale up from 0 failed", functionTarget.Function.Name))
			// send request completed event
			f2shub.F2SEventManager.Publish(eventmanager.Event{
				UID:  f2shub.F2SEventManager.GenerateUUID(),
				Data: req,
				Type: eventmanager.Event_FunctionInvokationEnded,
			})
			// send result to channel
			req.ResultChannel <- queue.F2SRequestResult{
				Result:  fmt.Sprintf("aborting function '%s' invocation because scale up from 0 failed", functionTarget.Function.Name),
				Success: false,
				UID:     req.UID,
			}
			return
		}
	}

	pod, err := functionTarget.ServeRequest(req)
	if err != nil {
		logging.Error(err)
		return
	}

	// start time measurement
	start := time.Now()

	// invoke function
	url := fmt.Sprintf("http://%s:%v%s", string(pod.Address.IP), functionTarget.Function.Target.Port, functionTarget.Function.Target.Endpoint)
	url = fmt.Sprintf("http://192.168.2.40:32343%s", functionTarget.Function.Target.Endpoint)
	result, err := httpGet(url)
	logging.Error(err)

	// measure time elapsed
	elapsed := time.Since(start).Milliseconds()
	elapsedPerInflight := int(time.Since(start).Milliseconds()) / len(pod.InflightRequests)
	logging.Info("Function execution time: %s\n", fmt.Sprintf("%vms", elapsed))
	logging.Info("Function execution time per inflight request: %sms\n", fmt.Sprintf("%v", elapsedPerInflight))

	// send result to channel
	req.ResultChannel <- queue.F2SRequestResult{
		Result:                     result,
		Success:                    true,
		UID:                        req.UID,
		Duration:                   float64(elapsed),
		DurationPerInflightRequest: float64(elapsedPerInflight),
	}

	// send request completed event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:  req.UID,
		Data: req,
		Type: eventmanager.Event_FunctionInvokationEnded,
	})
}
