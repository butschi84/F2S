package dispatcher

import (
	"butschi84/f2s/state/eventmanager"
	"butschi84/f2s/state/queue"
	"fmt"
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

	// search function target
	functionTarget, _ := f2shub.F2STargets.GetFunctionTargetByEndpoint(req.Path)
	pod := functionTarget.ServeRequest(req)

	// send 'function invoked' event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:      f2shub.F2SEventManager.GenerateUUID(),
		Data:     req.Path,
		Function: functionTarget.Function,
		Type:     eventmanager.Event_FunctionInvoked,
	})

	// invoke
	url := fmt.Sprintf("http://%s:%v%s", string(pod.Address.IP), functionTarget.Function.Target.Port, functionTarget.Function.Target.Endpoint)
	url = fmt.Sprintf("http://192.168.2.40:32343%s", functionTarget.Function.Target.Endpoint)
	result, _ := httpGet(url)

	// send result to channel
	req.ResultChannel <- result

	// send request completed event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:      f2shub.F2SEventManager.GenerateUUID(),
		Data:     req,
		Function: functionTarget.Function,
		Type:     eventmanager.Event_FunctionInvokationEnded,
	})
}
