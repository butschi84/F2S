package dispatcher

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/queue"

	corev1 "k8s.io/api/core/v1"
)

type FunctionServingPod struct {
	Address          corev1.EndpointAddress
	InflightRequests []queue.F2SRequest
}

type FunctionTarget struct {
	Function    typesV1alpha1.Function
	ServingPods []FunctionServingPod
}
type IFunctionTarget interface {
	ServeRequest(queue.F2SRequest)
}
