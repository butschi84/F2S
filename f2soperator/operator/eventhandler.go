package operator

import (
	"butschi84/f2s/eventmanager"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	kubernetesservice "butschi84/f2s/services/kubernetes"
)

func handleEvent(event eventmanager.Event) {
	logging.Info("processing event", fmt.Sprintf("'%s'", string(event.Type)))

	switch event.Type {
	case eventmanager.Event_ConfigurationChanged:
		if master {
			logging.Info("configuration has changed. rebalance immediately")
			Rebalance()
		}
	}
}

// will be called when f2sfunction / crd changes in k8s
func OnF2SFunctionChanged(obj interface{}) {
	logging.Info("F2S Functions Changes. Reloading Config...")

	logging.Info("read all F2SFunctions from K8S")
	functions, err := kubernetesservice.GetF2SFunctions()
	if err != nil {
		logging.Info("Failed to read f2s functions")
		return
	}

	// update active configuration
	F2SHub.F2SConfiguration.Functions = functions
	logging.Info("number of functions:", string(len(F2SHub.F2SConfiguration.Functions.Items)))

	// send config change event
	F2SHub.F2SEventManager.Publish(eventmanager.Event{
		UID:  F2SHub.F2SEventManager.GenerateUUID(),
		Data: "F2SFunctions Changed in K8S",
		Type: eventmanager.Event_ConfigurationChanged,
	})
}

func OnF2SEndpointsChanged(obj interface{}) {
	logging.Info("Endpoints have changed")

	// try parse endpoint obj
	d := &corev1.Endpoints{}
	err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.(*unstructured.Unstructured).UnstructuredContent(), d)
	if err != nil {
		logging.Error(fmt.Errorf("could not convert event to endpoint"))
		logging.Error(err)
		return
	}
	logging.Info(fmt.Sprintf("changed endpoint %s (%s)", d.Name, d.UID))
}
