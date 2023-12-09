package dispatcherstate

import (
	"butschi84/f2s/services/logger"
)

var logging *logger.F2SLogger

var functionTargets F2SDispatcherHub

func init() {
	// initialize logging
	logging = logger.Initialize("functiontargets")
}

func Initialize() *F2SDispatcherHub {
	functionTargets = F2SDispatcherHub{
		DispatcherFunctions: make([]F2SDispatcherFunction, 0),
	}
	return &functionTargets
}
