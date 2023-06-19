package kubernetesservice

import (
	"butschi84/f2s/logger"
	"log"
)

var logging *log.Logger

func init() {
	// initialize logging
	logging = logger.Initialize("kubernetesservice")
}

func int32Ptr(i int32) *int32 {
	return &i
}
