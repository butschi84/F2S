package dispatcher

import (
	"butschi84/f2s/hub"
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/queue"
	"sync"
)

type F2SDispatcher struct {
	Hub       *hub.F2SHub
	WaitGroup *sync.WaitGroup

	FunctionTargets []FunctionTarget
}
type IF2SDispatcher interface {
	HandleRequests()
}

type FunctionServingPod struct {
	IP               string
	InflightRequests []queue.F2SRequest
}

type FunctionTarget struct {
	Function    typesV1alpha1.Function
	ServingPods []FunctionServingPod
}
type IFunctionTarget interface {
	ServeRequest(queue.F2SRequest)
}
