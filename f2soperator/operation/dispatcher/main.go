package dispatcher

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"butschi84/f2s/state/queue"
	"fmt"
	"sync"

	kubernetesservice "butschi84/f2s/services/kubernetes"
)

var logging logger.F2SLogger
var f2shub *hub.F2SHub
var functionTargets []FunctionTarget

func Initialize(hub *hub.F2SHub, wg *sync.WaitGroup) {
	defer wg.Done()

	// consume variables
	f2shub = hub

	// initialize logging
	logging = logger.Initialize("dispatcher")

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	hub.F2SEventManager.Subscribe(handleEvents)

	// subscribe to new requests
	logging.Info("subscribing to incoming requests")
	hub.F2SQueue.Subscribe(handleRequests)

	// initialize f2stargets
	functionTargets = make([]FunctionTarget, len(hub.F2SConfiguration.Functions.Items))
	for i, function := range hub.F2SConfiguration.Functions.Items {
		// get endpoint IP's for function
		endpoint, _ := kubernetesservice.GetEndpointWithName(function.Name)

		// prepare target obj
		target := FunctionTarget{
			Function:    function,
			ServingPods: make([]FunctionServingPod, len(endpoint.Subsets[0].Addresses)),
		}

		for i, address := range endpoint.Subsets[0].Addresses {
			fsp := FunctionServingPod{
				Address:          address,
				InflightRequests: make([]queue.F2SRequest, 1),
			}
			target.ServingPods[i] = fsp
		}

		logging.Info(fmt.Sprintf("function %s has %v endpoints", function.Name, len(target.ServingPods)))

		functionTargets[i] = target
	}
}
