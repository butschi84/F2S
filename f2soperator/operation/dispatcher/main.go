package dispatcher

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"sync"
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
}
