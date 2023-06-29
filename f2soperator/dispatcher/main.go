package dispatcher

import (
	"butschi84/f2s/configuration"
	"butschi84/f2s/services/logger"
	"sync"
)

var logging logger.F2SLogger

type F2SDispatcher struct {
	Config    *configuration.F2SConfiguration
	WaitGroup *sync.WaitGroup
}

type IF2SDispatcher interface {
	HandleRequests()
}

func init() {
	// initialize logging
	logging = logger.Initialize("dispatcher")
}

// main dispatcher function
func (F2SDispatcher *F2SDispatcher) HandleRequests() {
	defer F2SDispatcher.WaitGroup.Done()

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	F2SDispatcher.Config.EventManager.Subscribe(handleEvents)
}
