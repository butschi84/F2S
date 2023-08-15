package operator

import (
	"butschi84/f2s/state/eventmanager"
	"fmt"
	"time"

	v1alpha1types "butschi84/f2s/state/configuration/api/types/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	kubernetesservice "butschi84/f2s/services/kubernetes"
)

func handleEvent(event eventmanager.Event) {
	logging.Info("processing event", fmt.Sprintf("'%s'", string(event.Type)))

	switch event.Type {
	case eventmanager.Event_ConfigurationChanged:
		if f2shub.F2SOperatorState.IsMaster {
			logging.Info("configuration has changed. rebalance immediately")
			Rebalance()
		}
	case eventmanager.Event_FunctionInvoked:
		logging.Info("function invoked. checking minimum availability")
		function := event.Data.(v1alpha1types.PrettyFunction)
		checkMinimumAvailability(&function)
	}
}

// make sure that there is at least a scale of 1 replica available
func checkMinimumAvailability(function *v1alpha1types.PrettyFunction) {
	logging.Info("[checkMinimumAvailability] checking minimum availability")
	target, err := f2shub.F2SDispatcherHub.GetDispatcherFunctionByName(function.Name)
	if err != nil {
		logging.Error(fmt.Errorf("[checkMinimumAvailability] could not get target for function %s. %s", function.Name, err.Error()))
		return
	}
	if len(target.ServingPods) == 0 {
		logging.Info(fmt.Sprintf("[checkMinimumAvailability] scaling up deployment %s to 1 replica", function.Name))
		kubernetesservice.ScaleDeployment(function.Name, 1)
		target.LastScaling = time.Now()

		// send 'function scaled' event
		f2shub.F2SEventManager.Publish(eventmanager.Event{
			UID:         f2shub.F2SEventManager.GenerateUUID(),
			Data:        function,
			Type:        eventmanager.Event_FunctionScaled,
			Description: fmt.Sprintf("F2SFunction %s scaled from %v to %v", function.Name, 0, 1),
		})
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
	f2shub.F2SConfiguration.Functions = functions
	logging.Info(fmt.Sprintf("number of functions: %d", len(f2shub.F2SConfiguration.Functions.Items)))

	// send config change event
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:         f2shub.F2SEventManager.GenerateUUID(),
		Data:        "F2SFunctions Changed in K8S",
		Type:        eventmanager.Event_ConfigurationChanged,
		Description: fmt.Sprintf("functions have changed in k8s"),
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

	// send change event
	logging.Info(fmt.Sprintf("send event for changed endpoint %s (%s)", d.Name, d.UID))
	f2shub.F2SEventManager.Publish(eventmanager.Event{
		UID:         f2shub.F2SEventManager.GenerateUUID(),
		Data:        d,
		Type:        eventmanager.Event_EndpointsChanged,
		Description: fmt.Sprintf("endpoint %s has changed in k8s", d.Name),
	})

}
