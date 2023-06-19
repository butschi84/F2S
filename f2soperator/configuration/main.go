package configuration

import (
	typesV1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"
	"butschi84/f2s/logger"
	"butschi84/f2s/services/eventmanager"
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"log"
)

var logging *log.Logger

type F2SConfiguration struct {
	Functions    *typesV1alpha1.FunctionList
	EventManager *eventmanager.EventManager
}

var ActiveConfiguration F2SConfiguration

func init() {
	// initialize logging
	logging = logger.Initialize("configuration")

	logging.Println("getting all f2s functions using kubernetes service")
	functions, err := kubernetesservice.GetF2SFunctions()
	if err != nil {
		logging.Println("Failed to read f2s functions")
		return
	}

	eventManager := eventmanager.NewEventManager()
	eventManager.Start()

	logging.Println("initializing config")
	ActiveConfiguration = F2SConfiguration{
		Functions:    functions,
		EventManager: eventManager,
	}

	logging.Println("starting to watch f2sfunctions in k8s")
	go kubernetesservice.WatchF2SFunctions(OnF2SFunctionChanged)
}
