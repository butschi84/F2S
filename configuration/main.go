package configuration

import (
	typesV1alpha1 "butschi84/f2s/configuration/api/types/v1alpha1"
	"butschi84/f2s/logger"
	"log"
	"time"
)

var logging *log.Logger

type F2SConfiguration struct {
	Functions *typesV1alpha1.FunctionList
}

var ActiveConfiguration F2SConfiguration

func init() {
	// initialize logging
	logging = logger.Initialize("configuration")

	logging.Println("initializing config")
	ActiveConfiguration = F2SConfiguration{
		Functions: GetCRDs(),
	}

	// reload every 30 seconds
	go reloadConfig()
}

func reloadConfig() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Your logic here, executed every 30 seconds
			logging.Println("reloading active config")
			ActiveConfiguration = F2SConfiguration{
				Functions: GetCRDs(),
			}
		}
	}
}
