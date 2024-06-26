package hub

import (
	clusterstate "butschi84/f2s/state/cluster"
	"butschi84/f2s/state/configuration"
	"butschi84/f2s/state/dispatcherstate"
	"butschi84/f2s/state/eventmanager"
	"butschi84/f2s/state/operatorstate"
	"butschi84/f2s/state/queue"
)

type F2SHub struct {
	F2SConfiguration *configuration.F2SConfiguration
	F2SEventManager  *eventmanager.EventManager
	F2SQueue         *queue.F2SQueue
	F2SDispatcherHub *dispatcherstate.F2SDispatcherHub
	F2SOperatorState *operatorstate.F2SOperatorState
	F2SClusterState  *clusterstate.F2SClusterState
}
