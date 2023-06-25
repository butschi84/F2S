package metrics

import (
	"butschi84/f2s/services/eventmanager"
	"butschi84/f2s/services/prometheus"
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
		currentInflightRequests += 1

	case eventmanager.Event_FunctionInvokationEnded:
		duration := event.Data.(time.Duration)

		// increase metric 'total_completed_requests
		logging.Println(fmt.Sprintf("function %s was invoked. increasing counter 'metricTotalCompletedRequests'", event.Data))
		metricTotalCompletedRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Inc()

		logging.Println(fmt.Sprintf("function %s finished. recalculating capacity", event.Data))
		numbercontainers, err := prometheus.ReadPrometheusMetricValue(&F2SConfiguration, "f2sscaling_function_deployment_available_replicas", map[string]string{"functionname": event.Function.Name})
		if err != nil {
			logging.Println("Error when trying to read prometheus metric f2sscaling_function_deployment_available_replicas", err)
		} else {
			inflightRequestsPerContainer := float64(currentInflightRequests) / numbercontainers
			durationPerInflightRequest := float64(duration.Seconds()) / inflightRequestsPerContainer
			functionCapacityRequestsPerSecond := 1 / durationPerInflightRequest
			metricFunctionCapacity.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Observe(functionCapacityRequestsPerSecond)
		}

		// decrease metric 'active_requests
		logging.Println(fmt.Sprintf("function %s finished. descreasing auge metricactiveRequests", event.Data))
		currentInflightRequests -= 1
		metricActiveRequests.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Dec()

		// update request duration metric
		logging.Println("observing function duration:", float64(duration.Seconds()))
		metricRequestDuration.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Observe(float64(duration.Seconds()))

		// update last request completion metric
		logging.Println(fmt.Sprintf("function %s finished. set timestamp of metric lastRequestCompletion", event.Data))
		metricLastRequestCompletion.WithLabelValues(event.Function.Spec.Endpoint, string(event.Function.UID), event.Function.Name).Set(float64(time.Now().Unix()))

	default:
		logging.Println("no action defined for", event.Type)
	}
}
