package dispatcher

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
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
		req := event.Data.(queue.F2SRequestResult).Request
		functionTarget, _ := f2shub.F2STargets.GetFunctionTargetByEndpoint(req.Path)
		functionTarget.RemoveRequest(req)
		f2shub.F2SQueue.RequestDone(req)
	case eventmanager.Event_EndpointsChanged:
		logging.Info("endpoints have changed")
		reloadEndpoints()
	}
}

// reload endpoint information in state when endpoints in k8s change
func reloadEndpoints() {
	logging.Info("reloading endpoint information")

	// add functions and f2sfunctiontargets
	for _, function := range f2shub.F2SConfiguration.Functions.Items {
		// get endpoint IP's for function
		endpoint, _ := kubernetesservice.GetEndpointWithName(function.Name)

		// prepare target obj
		target := f2shub.F2STargets.AddDispatcherFunction(&function)

		if len(endpoint.Subsets) > 0 {
			for _, address := range endpoint.Subsets[0].Addresses {
				logging.Info(fmt.Sprintf("add servingpod %s for function %s", string(address.IP), function.Name))
				target.AddServingPod(address)
			}

			// remove surplus serving-pods
			for _, pod := range target.ServingPods {
				validEndpoint := false
				for _, address := range endpoint.Subsets[0].Addresses {
					if string(address.IP) == string(pod.Address.IP) {
						validEndpoint = true
					}
				}
				if !validEndpoint {
					logging.Info(fmt.Sprintf("remove servingpod %s for function %s", string(pod.Address.IP), function.Name))
					target.RemoveServingPod(pod.Address)
				}
			}
		} else {
			logging.Info(fmt.Sprintf("function %s is scaled to zero", function.Name))
			target.ServingPods = make([]f2sfunctiontargets.FunctionServingPod, 0)
		}

		logging.Info(fmt.Sprintf("function %s has %v endpoints", function.Name, len(target.ServingPods)))
	}

}

// wait until f2s has scaled up the deploment from 0
func waitForTargetPod(target *f2sfunctiontargets.F2SDispatcherFunctionTarget) error {
	// Create a context with a timeout of 'scaling_timeout'
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
