package hub

import (
	"butschi84/f2s/state/configuration"
	"butschi84/f2s/state/eventmanager"
)

type F2SHub struct {
	F2SConfiguration configuration.F2SConfiguration
	F2SEventManager  *eventmanager.EventManager
}
