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
		request := event.Data.(queue.F2SRequest)
		logging.Info(fmt.Sprintf("function %s was invoked. increasing counter 'metricTotalIncomingRequests'", request.Function.Name))
		metricTotalIncomingRequests.WithLabelValues(request.Function.Spec.Endpoint, string(request.Function.UID), request.Function.Name, request.F2SUser.Username).Inc()

	case eventmanager.Event_FunctionScaled:
		function := event.Data.(v1alpha1types.PrettyFunction)
		logging.Info(fmt.Sprintf("function %s has just been scaled", function.Name))
		metricLastFunctionScaling.WithLabelValues(function.Spec.Endpoint, string(function.UID), function.Name).Set(float64(time.Now().Unix()))

	case eventmanager.Event_InflightRequestsChanged:
		request := event.Data.(queue.F2SRequest)
		logging.Info(fmt.Sprintf("number of inflight requests changed for function: %s", request.Function.Name))
		for _, dispatcherFunction := range f2shub.F2SDispatcherHub.DispatcherFunctions {
			inflightRequests := 0
			for _, servingPod := range dispatcherFunction.ServingPods {
				inflightRequests += len(servingPod.InflightRequests)
			}
			metricActiveRequests.WithLabelValues(dispatcherFunction.Function.Spec.Endpoint, string(dispatcherFunction.Function.UID), dispatcherFunction.Function.Name).Set(float64(inflightRequests))
		}
	case eventmanager.Event_FunctionInvokationEnded:
		result := event.Data.(queue.F2SRequestResult)
		functionTarget, err := f2shub.F2SDispatcherHub.GetDispatcherFunctionByEndpoint(result.Request.Path)
		if err != nil {
			logging.Error(fmt.Errorf("[%s] cannot calculate metrics for request because: %s", result.Request.UID, err.Error()))
		}

		if result.Success {
			// increase metric 'total_completed_requests
			logging.Info(fmt.Sprintf("function %s invokation ended. increasing counter 'metricTotalCompletedRequests'", functionTarget.Function.Name))
			metricTotalCompletedRequests.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Inc()

			// increase metric 'total request duration'
			logging.Info(fmt.Sprintf("function %s finished. increase metric 'total request duration' by %vms", functionTarget.Function.Name, result.Duration))
			metricTotalRequestDuration.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name, result.Request.F2SUser.Username).Add(result.Duration)

			// recaulculate capacity
			logging.Info(fmt.Sprintf("function %s finished. recalculating capacity", functionTarget.Function.Name))
			functionCapacityRequestsPerSecond := 1000 / result.DurationPerInflightRequest
			metricFunctionCapacity.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Observe(functionCapacityRequestsPerSecond)

			// update request duration metric
			logging.Info("observing function duration:", fmt.Sprintf("%v", float64(result.Duration)))
			metricRequestDuration.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Observe(result.Duration)
		} else {
			// increase metric 'total_failed_requests
			logging.Info(fmt.Sprintf("function %s invokation ended. increasing counter 'metricTotalFailedRequests'", functionTarget.Function.Name))
			metricTotalFailedRequests.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Inc()
		}

		// update last request completion metric
		logging.Info(fmt.Sprintf("function %s invokation ended. set timestamp of metric lastRequestCompletion", functionTarget.Function.Name))
		metricLastRequestCompletion.WithLabelValues(functionTarget.Function.Spec.Endpoint, string(functionTarget.Function.UID), functionTarget.Function.Name).Set(float64(time.Now().Unix()))

	default:
		logging.Info("no action defined for", fmt.Sprintf("%s", event.Type))
	}
}
