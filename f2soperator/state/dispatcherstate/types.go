package dispatcherstate

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"butschi84/f2s/state/queue"
	"time"

	corev1 "k8s.io/api/core/v1"
)

type DispatcherFunctionTarget struct {
	Address          corev1.EndpointAddress
	InflightRequests []queue.F2SRequest
}

type F2SDispatcherFunction struct {
	Function    typesV1alpha1.Function
	ServingPods []DispatcherFunctionTarget

	LastScaling time.Time
}
type IF2SDispatcherFunctionTarget interface {
	ServeRequest(queue.F2SRequest)
	AddServingPod(corev1.EndpointAddress)
	RemoveServingPod(corev1.EndpointAddress)
}
type F2SDispatcherHub struct {
	DispatcherFunctions []F2SDispatcherFunction
}

// ####################################
// Pretty Format for API / Frontend
// ####################################

type F2SDispatcherPrettyHub struct {
	DispatcherFunctions []F2SDispatcherPrettyFunction `json:"dispatcher_functions"`
}
type F2SDispatcherPrettyFunction struct {
	UID         string                           `json:"uid"`
	Name        string                           `json:"name"`
	Endpoint    string                           `json:"endpoint"`
	ServingPods []F2SDispatcherPrettyFunctionPod `json:"endpoints"`
}
type F2SDispatcherPrettyFunctionPod struct {
	IPAddress        string `json:"ip_address"`
	InflightRequests int    `json:"inflight_requests"`
}

func (f *F2SDispatcherFunction) Pretty() F2SDispatcherPrettyFunction {
	result := F2SDispatcherPrettyFunction{
		UID:         string(f.Function.UID),
		Name:        f.Function.Name,
		Endpoint:    f.Function.Spec.Endpoint,
		ServingPods: make([]F2SDispatcherPrettyFunctionPod, 0),
	}
	for _, p := range f.ServingPods {
		result.ServingPods = append(result.ServingPods, p.Pretty())
	}
	return result
}
func (fsp *DispatcherFunctionTarget) Pretty() F2SDispatcherPrettyFunctionPod {
	return F2SDispatcherPrettyFunctionPod{
		IPAddress:        fsp.Address.IP,
		InflightRequests: len(fsp.InflightRequests),
	}
}
func (hub *F2SDispatcherHub) Pretty() F2SDispatcherPrettyHub {
	result := F2SDispatcherPrettyHub{
		DispatcherFunctions: make([]F2SDispatcherPrettyFunction, 0),
	}
	for _, f := range hub.DispatcherFunctions {
		result.DispatcherFunctions = append(result.DispatcherFunctions, f.Pretty())
	}
	return result
}
