package f2sfunctiontargets

import (
	"butschi84/f2s/services/logger"
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/queue"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

var logging logger.F2SLogger

type FunctionServingPod struct {
	Address          corev1.EndpointAddress
	InflightRequests []queue.F2SRequest
}

type F2SFunctionTarget struct {
	Function    typesV1alpha1.Function
	ServingPods []FunctionServingPod
}
type IF2SFunctionTarget interface {
	ServeRequest(queue.F2SRequest)
}
type F2SDispatcherFunction struct {
	Targets []F2SFunctionTarget
}

var functionTargets F2SDispatcherFunction

func init() {
	// initialize logging
	logging = logger.Initialize("functiontargets")
}

func Initialize() *F2SDispatcherFunction {
	functionTargets = F2SDispatcherFunction{
		Targets: []F2SFunctionTarget{},
	}
	return &functionTargets
}

// put new incoming requests into queue
func (target *F2SFunctionTarget) ServeRequest(request queue.F2SRequest) (*FunctionServingPod, error) {
	// error prevention
	if len(target.ServingPods) == 0 {
		return &FunctionServingPod{}, fmt.Errorf("Cannot serve request %s. function %s has 0 pods available", request.UID, request.Path)
	}

	// add inflight request to first pod in array
	target.ServingPods[0].InflightRequests = append(target.ServingPods[0].InflightRequests, request)

	// shift pods, so next invocation is served by another pod (1,2,3 => 2,3,1)
	if len(target.ServingPods) > 1 {
		logging.Info("rotating serving pods")

		firstPod := target.ServingPods[0]
		s2 := make([]FunctionServingPod, len(target.ServingPods))
		copy(s2, target.ServingPods[1:])
		s2[len(target.ServingPods)-1] = firstPod
		copy(target.ServingPods, s2)
	}

	return &target.ServingPods[len(target.ServingPods)-1], nil
}

func (target *F2SFunctionTarget) RemoveRequest(request queue.F2SRequest) {
	for i, pod := range target.ServingPods {
		for x, req := range pod.InflightRequests {
			if req.UID == request.UID {
				target.ServingPods[i].InflightRequests = append(
					target.ServingPods[i].InflightRequests[:x],
					target.ServingPods[i].InflightRequests[x+1:]...,
				)
				return
			}
		}
	}
}

// get target by path
func (target *F2SDispatcherFunction) GetFunctionTargetByEndpoint(endpoint string) (*F2SFunctionTarget, error) {
	logging.Info(fmt.Sprintf("get functiontarget for path %s", endpoint))
	for _, target := range target.Targets {
		if target.Function.Spec.Endpoint == endpoint {
			return &target, nil
		}
	}

	return &F2SFunctionTarget{}, fmt.Errorf("functiontarget not found")
}

// get target by path
func (target *F2SDispatcherFunction) GetFunctionTargetByFunctionName(name string) (*F2SFunctionTarget, error) {
	logging.Info(fmt.Sprintf("get functiontarget for function name %s", name))
	for _, target := range target.Targets {
		if target.Function.Name == name {
			return &target, nil
		}
	}

	return &F2SFunctionTarget{}, fmt.Errorf("functiontarget not found")
}
