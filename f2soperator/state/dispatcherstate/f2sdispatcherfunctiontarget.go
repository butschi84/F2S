package dispatcherstate

import (
	"butschi84/f2s/state/queue"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// put new incoming requests into queue
func (dispatcherFunction *F2SDispatcherFunction) ServeRequest(request *queue.F2SRequest) (*DispatcherFunctionTarget, error) {
	// error prevention
	if len(dispatcherFunction.ServingPods) == 0 {
		return &DispatcherFunctionTarget{}, fmt.Errorf("Cannot serve request %s. function %s has 0 pods available", request.UID, request.Path)
	}

	// add inflight request to first pod in array
	dispatcherFunction.ServingPods[0].InflightRequests = append(dispatcherFunction.ServingPods[0].InflightRequests, *request)

	// shift pods, so next invocation is served by another pod (1,2,3 => 2,3,1)
	if len(dispatcherFunction.ServingPods) > 1 {
		logging.Info("rotating serving pods")

		firstPod := dispatcherFunction.ServingPods[0]
		s2 := make([]DispatcherFunctionTarget, len(dispatcherFunction.ServingPods))
		copy(s2, dispatcherFunction.ServingPods[1:])
		s2[len(dispatcherFunction.ServingPods)-1] = firstPod
		copy(dispatcherFunction.ServingPods, s2)
	}
	return &dispatcherFunction.ServingPods[len(dispatcherFunction.ServingPods)-1], nil
}

func (dispatcherFunction *F2SDispatcherFunction) GetTotalInflightRequests() int {
	total := 0
	for _, pod := range dispatcherFunction.ServingPods {
		total += len(pod.InflightRequests)
	}
	return total
}

func (dispatcherFunction *F2SDispatcherFunction) RemoveRequest(request queue.F2SRequest) {
	for i, pod := range dispatcherFunction.ServingPods {
		for x, req := range pod.InflightRequests {
			if req.UID == request.UID {
				dispatcherFunction.ServingPods[i].InflightRequests = append(
					dispatcherFunction.ServingPods[i].InflightRequests[:x],
					dispatcherFunction.ServingPods[i].InflightRequests[x+1:]...,
				)
				return
			}
		}
	}
}

func (dispatcherFunction *F2SDispatcherFunction) SetLastScaling() {
	dispatcherFunction.LastScaling = time.Now()
}

func (dispatcherFunction *F2SDispatcherFunction) AddServingPod(address corev1.EndpointAddress, podUID string, podName string) {
	logging.Debug(fmt.Sprintf("add servingpod %s to dispatchertarget %s", string(address.IP), dispatcherFunction.Function.Name))
	for _, pod := range dispatcherFunction.ServingPods {
		if string(pod.Address.IP) == string(address.IP) {
			logging.Debug(fmt.Sprintf("dispatchertarget %s already has servingpod %s", dispatcherFunction.Function.Name, string(address.IP)))
			return
		}
	}
	dispatcherFunction.ServingPods = append(dispatcherFunction.ServingPods, DispatcherFunctionTarget{
		InflightRequests: make([]queue.F2SRequest, 0),
		Address:          address,
		UID:              podUID,
		Name:             podName,
	})
}

func (dispatcherFunction *F2SDispatcherFunction) RemoveServingPod(address corev1.EndpointAddress) {
	for i, pod := range dispatcherFunction.ServingPods {
		if string(pod.Address.IP) == string(address.IP) {
			// remove from array
			dispatcherFunction.ServingPods = append(dispatcherFunction.ServingPods[:i], dispatcherFunction.ServingPods[i+1:]...)
		}
	}
}
