package configuration

import (
	"butschi84/f2s/services/eventmanager"
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// will be called when f2sfunction / crd changes in k8s
func OnF2SFunctionChanged(obj interface{}) {
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
		Type: eventmanager.Event_ConfigurationChanged,
	})
}

func OnF2SEndpointsChanged(obj interface{}) {
	logging.Println("Endpoints have changed")

	// try parse endpoint obj
	d := &corev1.Endpoints{}
	err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.(*unstructured.Unstructured).UnstructuredContent(), d)
	if err != nil {
		logging.Println("could not convert event to endpoint")
		logging.Print(err)
		return
	}
	logging.Println(fmt.Sprintf("changed endpoint %s (%s)", d.Name, d.UID))
}
