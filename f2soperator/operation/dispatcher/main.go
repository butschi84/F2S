package dispatcher

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"fmt"
	"sync"
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

	reloadEndpoints()
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
