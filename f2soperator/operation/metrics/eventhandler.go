package metrics

import (
	v1alpha1types "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/eventmanager"
	"butschi84/f2s/state/queue"
	"fmt"
	"time"
)

func handleEvent(event eventmanager.Event) {
	switch event.Type {
	// function invoked
	case eventmanager.Event_FunctionInvoked:
		// increase metric 'total_incoming_requests
		logging.Info(fmt.Sprintf("function %s was invoked. increasing counter 'metricTotalIncomingRequests'", event.Data))
		prettyFunction := event.Data.(v1alpha1types.PrettyFunction)
		metricTotalIncomingRequests.WithLabelValues(prettyFunction.Spec.Endpoint, string(prettyFunction.UID), prettyFunction.Name).Inc()

		// increase metric 'active_requests
		logging.Info(fmt.Sprintf("function %s was invoked. increasing counter 'metricactiveRequests'", event.Data))
		metricActiveRequests.WithLabelValues(prettyFunction.Spec.Endpoint, string(prettyFunction.UID), prettyFunction.Name).Inc()
		currentInflightRequests += 1

	case eventmanager.Event_FunctionInvokationEnded:
		result := event.Data.(queue.F2SRequestResult)
		functionTarget, err := f2shub.F2STargets.GetFunctionTargetByEndpoint(result.Request.Path)
		if err != nil {
			logging.Error(fmt.Errorf("cannot calculate metrics for request %s", result.Request.UID))
			logging.Error(err)
		}

		// increase metric 'total_completed_requests
		logging.Info(fmt.Sprintf("function %s was invoked. increasing counter 'metricTotalCompletedRequests'", event.Data))
		metricTotalCompletedRequests.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Inc()

		logging.Info(fmt.Sprintf("function %s finished. recalculating capacity", functionTarget.Function.Name))
		functionCapacityRequestsPerSecond := 1000 / result.DurationPerInflightRequest
		metricFunctionCapacity.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Observe(functionCapacityRequestsPerSecond)

		// decrease metric 'active_requests
		logging.Info(fmt.Sprintf("function %s finished. descreasing auge metricactiveRequests", event.Data))
		currentInflightRequests -= 1
		metricActiveRequests.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Dec()

		// update request duration metric
		logging.Info("observing function duration:", fmt.Sprintf("%v", float64(result.Duration)))
		metricRequestDuration.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Observe(result.Duration)

		// update last request completion metric
		logging.Info(fmt.Sprintf("function %s finished. set timestamp of metric lastRequestCompletion", event.Data))
		metricLastRequestCompletion.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Set(float64(time.Now().Unix()))

	default:
		logging.Info("no action defined for", fmt.Sprintf("%s", event.Type))
	}
}
