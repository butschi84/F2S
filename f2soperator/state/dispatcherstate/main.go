package dispatcherstate

import (
	"butschi84/f2s/services/logger"

	"golang.org/x/exp/slog"
)

var logging *slog.Logger

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
