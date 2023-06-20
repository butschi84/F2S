package metrics

import (
	"butschi84/f2s/services/eventmanager"
	"fmt"
	"time"
)

func handleEvent(event eventmanager.Event) {
	switch event.Type {
	// function invoked
	case eventmanager.Event_FunctionInvoked:
		// increase metric 'total_requests
		logging.Println(fmt.Sprintf("function %s was invoked. increasing counter 'metricTotalRequests'", event.Data))
		metricTotalRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Inc()
		// increase metric 'active_requests
		logging.Println(fmt.Sprintf("function %s was invoked. increasing counter 'metricactiveRequests'", event.Data))
		metricActiveRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Inc()

	case eventmanager.Event_FunctionInvokationEnded:
		// decrease metric 'active_requests
		logging.Println(fmt.Sprintf("function %s finished. descreasing auge metricactiveRequests", event.Data))
		metricActiveRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Dec()

		// update request duration metric
		duration := event.Data.(time.Duration)
		durationInSeconds := float64(duration) / 1000.0
		metricRequestDuration.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Observe(durationInSeconds)

	default:
		logging.Println("no action defined for", event.Type)
	}
}
