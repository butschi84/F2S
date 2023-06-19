package configuration

import (
	"butschi84/f2s/services/eventmanager"
	kubernetesservice "butschi84/f2s/services/kubernetes"
)

// will be called when f2sfunction / crd changes in k8s
func OnF2SFunctionChanged() {
	logging.Println("F2S Functions Changes. Reloading Config...")

	logging.Println("read all F2SFunctions from K8S")
	functions, err := kubernetesservice.GetF2SFunctions()
	if err != nil {
		logging.Println("Failed to read f2s functions")
		return
	}

	// update active configuration
	ActiveConfiguration.Functions = functions
	logging.Println("number of functions:", len(ActiveConfiguration.Functions.Items))

	// send config change event
	ActiveConfiguration.EventManager.Publish(eventmanager.Event{
		Data: "F2SFunctions Changed in K8S",
	})
}
