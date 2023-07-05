package f2sfunctiontargets

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"fmt"
)

func (dispatcherfunction *F2SDispatcherFunction) AddDispatcherFunction(function *typesV1alpha1.Function) *F2SDispatcherFunctionTarget {
	for i, fuc := range dispatcherfunction.Targets {
		if fuc.Function.UID == function.UID {
			return &dispatcherfunction.Targets[i]
		}
	}

	dispatcherfunction.Targets = append(dispatcherfunction.Targets, F2SDispatcherFunctionTarget{
		Function:    *function,
		ServingPods: make([]FunctionServingPod, 0),
	})
	return &dispatcherfunction.Targets[len(dispatcherfunction.Targets)-1]
}

// get target by path
func (target *F2SDispatcherFunction) GetFunctionTargetByEndpoint(endpoint string) (*F2SDispatcherFunctionTarget, error) {
	logging.Debug(fmt.Sprintf("get functiontarget for path %s", endpoint))
	for i, t := range target.Targets {
		if t.Function.Spec.Endpoint == endpoint {
			return &target.Targets[i], nil
		}
	}

	return &F2SDispatcherFunctionTarget{}, fmt.Errorf("functiontarget not found")
}

// get target by path
func (target *F2SDispatcherFunction) GetFunctionTargetByFunctionName(name string) (*F2SDispatcherFunctionTarget, error) {
	logging.Debug(fmt.Sprintf("get functiontarget for function name %s", name))
	for i, t := range target.Targets {
		if t.Function.Name == name {
			return &target.Targets[i], nil
		}
	}

	return &F2SDispatcherFunctionTarget{}, fmt.Errorf("functiontarget not found")
}
