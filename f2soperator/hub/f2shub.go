package hub

import (
	"butschi84/f2s/state/configuration"
	"butschi84/f2s/state/eventmanager"
	f2sfunctiontargets "butschi84/f2s/state/functiontargets"
	"butschi84/f2s/state/queue"
)

type F2SHub struct {
	F2SConfiguration configuration.F2SConfiguration
	F2SEventManager  *eventmanager.EventManager
	F2SQueue         *queue.F2SQueue
	F2STargets       *f2sfunctiontargets.F2SFunctionTargets
}
