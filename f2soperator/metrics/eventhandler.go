package metrics

import (
	"butschi84/f2s/services/eventmanager"
	"fmt"
)

func handleEvent(event eventmanager.Event) {
	switch event.Type {
	// function invoked
	case eventmanager.Event_FunctionInvoked:
		logging.Println(fmt.Sprintf("function %s was invoked. increasing counter 'metricTotalRequests'", event.Data))
		metricTotalRequests.Inc()
	default:
		logging.Println("no action defined for", event.Type)
	}
}
