package f2sfunctiontargets

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"fmt"
)

// add a dispatcherfunction if not already present
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

	return &F2SDispatcherFunctionTarget{}, fmt.Errorf("functiontarget not found for endpoint: %s", endpoint)
}

// get target by path
func (target *F2SDispatcherFunction) GetFunctionTargetByFunctionName(functionName string) (*F2SDispatcherFunctionTarget, error) {
	logging.Debug(fmt.Sprintf("get functiontarget for function name %s", functionName))
	for i, t := range target.Targets {
		if t.Function.Name == functionName {
			return &target.Targets[i], nil
		}
	}

	return &F2SDispatcherFunctionTarget{}, fmt.Errorf("functiontarget not found for function: %s", functionName)
}
