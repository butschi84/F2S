package configuration

import (
	kubernetesservice "butschi84/f2s/services/kubernetes"
)

// will be called when f2sfunction / crd changes in k8s
func OnF2SFunctionChanged() {
	logging.Println("F2S Functions Changes. Reloading Config...")

	functions, err := kubernetesservice.GetF2SFunctions()
	if err != nil {
		logging.Println("Failed to read f2s functions")
		return
	}

	ActiveConfiguration.Functions = functions
	logging.Println("number of functions:", len(ActiveConfiguration.Functions.Items))
}
