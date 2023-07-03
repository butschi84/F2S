package dispatcher

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	f2sfunctiontargets "butschi84/f2s/state/functiontargets"
	"butschi84/f2s/state/queue"
	"fmt"
	"sync"

	kubernetesservice "butschi84/f2s/services/kubernetes"
)

var logging logger.F2SLogger
var f2shub *hub.F2SHub

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
	hub.F2STargets.Targets = make([]f2sfunctiontargets.F2SFunctionTarget, len(hub.F2SConfiguration.Functions.Items))
	for i, function := range hub.F2SConfiguration.Functions.Items {
		// get endpoint IP's for function
		endpoint, _ := kubernetesservice.GetEndpointWithName(function.Name)

		// prepare target obj
		hub.F2STargets.Targets[i] = f2sfunctiontargets.F2SFunctionTarget{
			Function:    function,
			ServingPods: make([]f2sfunctiontargets.FunctionServingPod, len(endpoint.Subsets[0].Addresses)),
		}

		for x, address := range endpoint.Subsets[0].Addresses {
			hub.F2STargets.Targets[i].ServingPods[x] = f2sfunctiontargets.FunctionServingPod{
				Address:          address,
				InflightRequests: make([]queue.F2SRequest, 0),
			}
		}

		logging.Info(fmt.Sprintf("function %s has %v endpoints", function.Name, len(hub.F2STargets.Targets[i].ServingPods)))
	}
}

// debug output for dispatcher troubleshooting
func GetCurrentDispatcherData() string {
	output := ""
	output += "Dispatcher Data"
	output += "==============="
	// iterate functions
	for _, function := range f2shub.F2SConfiguration.Functions.Items {
		output += function.Name

		// get function target
		target, err := f2shub.F2STargets.GetFunctionTargetByFunctionName(function.Name)
		logging.Error(err)

		output += fmt.Sprintf("Endpoints: %d", len(target.ServingPods))
		for _, endpoint := range target.ServingPods {
			output += fmt.Sprintf("=> %s (inflight requests: %d)", string(endpoint.Address.IP), len(endpoint.InflightRequests))
		}
	}
	return output
}
