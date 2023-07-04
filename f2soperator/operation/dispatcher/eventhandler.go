package dispatcher

import (
	"butschi84/f2s/state/eventmanager"
	f2sfunctiontargets "butschi84/f2s/state/functiontargets"
	"butschi84/f2s/state/queue"
	"context"
	"fmt"
	"time"
)

// handle eventmanager events
func handleEvents(event eventmanager.Event) {
	logging.Info("processing event", fmt.Sprintf("'%s'", string(event.Type)))

	switch event.Type {
	case eventmanager.Event_FunctionInvokationEnded:
		logging.Info("removing request from targets")
		req := event.Data.(queue.F2SRequest)
		functionTarget, _ := f2shub.F2STargets.GetFunctionTargetByEndpoint(req.Path)
		functionTarget.RemoveRequest(req)
		f2shub.F2SQueue.RequestDone(req)
	}
}

// handle function invocations
func handleRequests(req queue.F2SRequest) {
	logging.Info("processing new function invocation request")

	// find function target
	functionTarget, err := f2shub.F2STargets.GetFunctionTargetByEndpoint(req.Path)
	if err != nil {
		logging.Error(fmt.Errorf("cannot serve request %s. function target not found for endpoint %s", req.UID, req.Path))
		logging.Error(err)
	}

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

	// send 'function invoked' event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:  f2shub.F2SEventManager.GenerateUUID(),
		Data: req.Path,
		Type: eventmanager.Event_FunctionInvoked,
	})

	// invoke function
	url := fmt.Sprintf("http://%s:%v%s", string(pod.Address.IP), functionTarget.Function.Target.Port, functionTarget.Function.Target.Endpoint)
	url = fmt.Sprintf("http://192.168.2.40:32343%s", functionTarget.Function.Target.Endpoint)
	result, err := httpGet(url)
	logging.Error(err)

	// send result to channel
	req.ResultChannel <- queue.F2SRequestResult{
		Result:  result,
		Success: true,
		UID:     req.UID,
	}

	// send request completed event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:  req.UID,
		Data: req,
		Type: eventmanager.Event_FunctionInvokationEnded,
	})
}

// wait until f2s has scaled up the deploment from 0
func waitForTargetPod(target *f2sfunctiontargets.F2SFunctionTarget) error {
	// Create a context with a timeout of 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resultChan := make(chan int)

	// Start the goroutine
	go func() {
		for len(target.ServingPods) < 1 {
			time.Sleep(1 * time.Second)
		}
		resultChan <- len(target.ServingPods)
	}()

	select {
	case result := <-resultChan:
		logging.Info(fmt.Sprintf("%d pods are available to serve the function '%s'", result, target.Function.Name))
		return nil
	case <-ctx.Done():
		logging.Warn(fmt.Sprintf("0 pods are available to serve the function '%s'!", target.Function.Name))
		return fmt.Errorf("0 pods are available to serve the function '%s'!", target.Function.Name)
	}
}
