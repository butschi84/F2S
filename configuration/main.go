package configuration

import (
	typesV1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"
	"butschi84/f2s/logger"
	kubernetesservice "butschi84/f2s/services/kubernetes"
	"log"
)

var logging *log.Logger

type F2SConfiguration struct {
	Functions *typesV1alpha1.FunctionList
}

var ActiveConfiguration F2SConfiguration

func init() {
	// initialize logging
	logging = logger.Initialize("configuration")

	functions, err := kubernetesservice.GetF2SFunctions()
	if err != nil {
		logging.Println("Failed to read f2s functions")
		return
	}

	logging.Println("initializing config")
	ActiveConfiguration = F2SConfiguration{
		Functions: functions,
	}

	go kubernetesservice.WatchF2SFunctions(OnF2SFunctionChanged)
}
