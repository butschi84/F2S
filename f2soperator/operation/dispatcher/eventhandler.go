package dispatcher

import (
	"butschi84/f2s/state/eventmanager"
	"butschi84/f2s/state/queue"
	"fmt"
)

// handle eventmanager events
func handleEvents(event eventmanager.Event) {
	logging.Info("processing event", fmt.Sprintf("'%s'", string(event.Type)))

}

// handle function invocations
func handleRequests(req queue.F2SRequest) {
	logging.Info("processing new function invocation request")

	// search function target
}

// get FunctionTarget for Request
func getFunctionTarget(req queue.F2SRequest) {
	// for _, traget := range
}
