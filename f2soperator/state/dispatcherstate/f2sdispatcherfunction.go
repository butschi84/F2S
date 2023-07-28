package dispatcherstate

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"
	"fmt"
)

// add a dispatcherfunction if not already present
func (dispatcherfunction *F2SDispatcherHub) AddDispatcherFunction(function *typesV1alpha1.Function) *F2SDispatcherFunction {
	for i, fuc := range dispatcherfunction.DispatcherFunctions {
		if fuc.Function.UID == function.UID {
			return &dispatcherfunction.DispatcherFunctions[i]
		}
	}

	dispatcherfunction.DispatcherFunctions = append(dispatcherfunction.DispatcherFunctions, F2SDispatcherFunction{
		Function:    *function,
		ServingPods: make([]DispatcherFunctionTarget, 0),
	})
	return &dispatcherfunction.DispatcherFunctions[len(dispatcherfunction.DispatcherFunctions)-1]
}

// get target by path
func (dispatcherhub *F2SDispatcherHub) GetFunctionTargetByEndpoint(endpoint string) (*F2SDispatcherFunction, error) {
	logging.Debug(fmt.Sprintf("get functiontarget for path %s", endpoint))
	for i, t := range dispatcherhub.DispatcherFunctions {
		if t.Function.Spec.Endpoint == endpoint {
			return &dispatcherhub.DispatcherFunctions[i], nil
		}
	}

	return &F2SDispatcherFunction{}, fmt.Errorf("functiontarget not found for endpoint: %s", endpoint)
}

// get target by path
func (dispatcherhub *F2SDispatcherHub) GetFunctionTargetByFunctionName(functionName string) (*F2SDispatcherFunction, error) {
	logging.Debug(fmt.Sprintf("get functiontarget for function name %s", functionName))
	for i, t := range dispatcherhub.DispatcherFunctions {
		if t.Function.Name == functionName {
			return &dispatcherhub.DispatcherFunctions[i], nil
		}
	}

	return &F2SDispatcherFunction{}, fmt.Errorf("functiontarget not found for function: %s", functionName)
}
