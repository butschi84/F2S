package dispatcher

import (
	"butschi84/f2s/hub"
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"sync"
)

// F2SRequest to invoke a function
type F2SRequest struct {
	UID    string
	Path   string
	Method string
}

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
	InflightRequests []F2SRequest
}

type FunctionTarget struct {
	Function    typesV1alpha1.Function
	ServingPods []FunctionServingPod
}
type IFunctionTarget interface {
	ServeRequest(F2SRequest)
}
