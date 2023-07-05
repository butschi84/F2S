package f2sfunctiontargets

import (
	"butschi84/f2s/services/logger"
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/queue"

	corev1 "k8s.io/api/core/v1"
)

var logging logger.F2SLogger

type FunctionServingPod struct {
	Address          corev1.EndpointAddress
	InflightRequests []queue.F2SRequest
}

type F2SDispatcherFunctionTarget struct {
	Function    typesV1alpha1.Function
	ServingPods []FunctionServingPod
}
type IF2SDispatcherFunctionTarget interface {
	ServeRequest(queue.F2SRequest)
	AddServingPod(corev1.EndpointAddress)
	RemoveServingPod(corev1.EndpointAddress)
}
type F2SDispatcherFunction struct {
	Targets []F2SDispatcherFunctionTarget
}
type IF2SDispatcherFunction struct {
	AddTarget (F2SDispatcherFunctionTarget)
}

var functionTargets F2SDispatcherFunction

func init() {
	// initialize logging
	logging = logger.Initialize("functiontargets")
}

func Initialize() *F2SDispatcherFunction {
	functionTargets = F2SDispatcherFunction{
		Targets: make([]F2SDispatcherFunctionTarget, 0),
	}
	return &functionTargets
}
