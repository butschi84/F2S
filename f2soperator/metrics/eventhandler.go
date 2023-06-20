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
		// increase metric 'total_incoming_requests
		logging.Println(fmt.Sprintf("function %s was invoked. increasing counter 'metricTotalIncomingRequests'", event.Data))
		metricTotalIncomingRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Inc()

		// increase metric 'active_requests
		logging.Println(fmt.Sprintf("function %s was invoked. increasing counter 'metricactiveRequests'", event.Data))
		metricActiveRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Inc()

	case eventmanager.Event_FunctionInvokationEnded:
		// increase metric 'total_completed_requests
		logging.Println(fmt.Sprintf("function %s was invoked. increasing counter 'metricTotalCompletedRequests'", event.Data))
		metricTotalCompletedRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Inc()

		// decrease metric 'active_requests
		logging.Println(fmt.Sprintf("function %s finished. descreasing auge metricactiveRequests", event.Data))
		metricActiveRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Dec()

		// update request duration metric
		duration := event.Data.(time.Duration)
		durationInSeconds := float64(duration) / 1000.0
		metricRequestDuration.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Observe(durationInSeconds)

		// update last request completion metric
		logging.Println(fmt.Sprintf("function %s finished. set timestamp of metric lastRequestCompletion", event.Data))
		metricLastRequestCompletion.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Set(float64(time.Now().Unix()))

	default:
		logging.Println("no action defined for", event.Type)
	}
}
