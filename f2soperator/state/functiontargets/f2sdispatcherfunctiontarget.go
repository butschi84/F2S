package f2sfunctiontargets

import (
	"butschi84/f2s/state/queue"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// put new incoming requests into queue
func (target *F2SDispatcherFunctionTarget) ServeRequest(request queue.F2SRequest) (*FunctionServingPod, error) {
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

func (target *F2SDispatcherFunctionTarget) RemoveRequest(request queue.F2SRequest) {
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

func (target *F2SDispatcherFunctionTarget) AddServingPod(address corev1.EndpointAddress) {
	logging.Debug(fmt.Sprintf("add servingpod %s to dispatchertarget %s", string(address.IP), target.Function.Name))
	for _, pod := range target.ServingPods {
		if string(pod.Address.IP) == string(address.IP) {
			logging.Debug(fmt.Sprintf("dispatchertarget %s already has servingpod %s", target.Function.Name, string(address.IP)))
			return
		}
	}
	target.ServingPods = append(target.ServingPods, FunctionServingPod{
		InflightRequests: make([]queue.F2SRequest, 0),
		Address:          address,
	})
}

func (target *F2SDispatcherFunctionTarget) RemoveServingPod(address corev1.EndpointAddress) {
	for i, pod := range target.ServingPods {
		if string(pod.Address.IP) == string(address.IP) {
			// remove from array
			target.ServingPods = append(target.ServingPods[:i], target.ServingPods[i+1:]...)
		}
	}
}
