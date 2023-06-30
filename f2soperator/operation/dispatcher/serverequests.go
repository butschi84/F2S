package dispatcher

import (
	"butschi84/f2s/state/queue"
)

// put new incoming requests into queue
func (target *FunctionTarget) ServeRequest(request queue.F2SRequest) {
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
