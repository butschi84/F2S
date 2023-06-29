package dispatcher

import (
	"butschi84/f2s/services/logger"
)

var logging logger.F2SLogger

func init() {
	// initialize logging
	logging = logger.Initialize("dispatcher")
}

// main dispatcher function
func (F2SDispatcher *F2SDispatcher) HandleRequests() {
	defer F2SDispatcher.WaitGroup.Done()

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	F2SDispatcher.Hub.F2SConfiguration.EventManager.Subscribe(handleEvents)
}
