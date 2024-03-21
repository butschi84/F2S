package dispatcher

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"butschi84/f2s/state/dispatcherstate"
	"butschi84/f2s/state/eventmanager"
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
		functionTarget, _ := f2shub.F2SDispatcherHub.GetDispatcherFunctionByEndpoint(req.Path)
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
		target := f2shub.F2SDispatcherHub.AddDispatcherFunction(&function)

		if len(endpoint.Subsets) > 0 {
			for i, address := range endpoint.Subsets[0].Addresses {
				logging.Info(fmt.Sprintf("add servingpod %s for function %s", string(address.IP), function.Name))
				podUID := endpoint.Subsets[0].Addresses[i].TargetRef.UID
				podName := endpoint.Subsets[0].Addresses[i].TargetRef.Name
				target.AddServingPod(address, string(podUID), podName)
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
			target.ServingPods = make([]dispatcherstate.DispatcherFunctionTarget, 0)
		}

		logging.Info(fmt.Sprintf("function %s has %v endpoints", function.Name, len(target.ServingPods)))
	}

}

// wait until f2s has scaled up the deploment from 0
func waitForTargetPod(dispatcherFunction *dispatcherstate.F2SDispatcherFunction) error {
	scalingTimeout := f2shub.F2SConfiguration.Config.F2S.Timeouts.ScalingTimeout

	// Create a context with a timeout of 'scaling_timeout'
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(scalingTimeout)*time.Millisecond)
	defer cancel()

	resultChan := make(chan int)

	// Start the goroutine
	go func() {
		for len(dispatcherFunction.ServingPods) < 1 {
			time.Sleep(1 * time.Second)
		}
		resultChan <- len(dispatcherFunction.ServingPods)
	}()

	select {
	case result := <-resultChan:
		logging.Info(fmt.Sprintf("%d pods are available to serve the function '%s'", result, dispatcherFunction.Function.Name))
		return nil
	case <-ctx.Done():
		logging.Warn(fmt.Sprintf("scaling_timeout: 0 pods are available to serve the function '%s' after %dms!", dispatcherFunction.Function.Name, scalingTimeout))
		return fmt.Errorf("scaling_timeout: 0 pods are available to serve the function '%s' after %dms!", dispatcherFunction.Function.Name, scalingTimeout)
	}
}
