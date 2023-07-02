package f2sfunctiontargets

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/queue"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

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
type F2SFunctionTargets struct {
	Targets []F2SFunctionTarget
}

func Initialize() *F2SFunctionTargets {
	return &F2SFunctionTargets{}
}

// put new incoming requests into queue
func (target *F2SFunctionTarget) ServeRequest(request queue.F2SRequest) {
	// add inflight request to first pod in array
	target.ServingPods[0].InflightRequests = append(target.ServingPods[0].InflightRequests, request)

	// shift pods, so next invocation is served by another pod (1,2,3 => 2,3,1)
	if len(target.ServingPods) > 1 {
		firstPod := target.ServingPods[0]
		s2 := make([]FunctionServingPod, len(target.ServingPods))
		copy(s2, target.ServingPods[1:])
		s2[len(target.ServingPods)-1] = firstPod
	}
}

// get target by path
func (target *F2SFunctionTargets) GetFunctionTargetByEndpoint(endpoint string) (F2SFunctionTarget, error) {
	for _, target := range target.Targets {
		if target.Function.Spec.Endpoint == endpoint {
			return target, nil
		}
	}

	return F2SFunctionTarget{}, fmt.Errorf("functiontarget not found")
}
