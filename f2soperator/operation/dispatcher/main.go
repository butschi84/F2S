package dispatcher

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"fmt"
	"time"
)

var logging logger.F2SLogger
var f2shub *hub.F2SHub

func Initialize(h *hub.F2SHub) {

	// consume variables
	f2shub = h

	// initialize logging
	logging = logger.Initialize("dispatcher")

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	f2shub.F2SEventManager.Subscribe(handleEvents)

	// subscribe to new requests
	logging.Info("subscribing to incoming requests")
	f2shub.F2SQueue.Subscribe(handleRequestsWithTimeout)

	reloadEndpoints()

	for {
		time.Sleep(time.Second)
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
		target, err := f2shub.F2SDispatcherHub.GetDispatcherFunctionByName(function.Name)
		logging.Error(err)

		output += fmt.Sprintf("Endpoints: %d", len(target.ServingPods))
		for _, endpoint := range target.ServingPods {
			output += fmt.Sprintf("=> %s (inflight requests: %d)", string(endpoint.Address.IP), len(endpoint.InflightRequests))
		}
	}
	return output
}
