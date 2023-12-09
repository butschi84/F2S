package v1alpha1

import (
	"butschi84/f2s/services/logger"
)

var logging *logger.F2SLogger

func init() {
	// initialize logging
	logging = logger.Initialize("v1alpha1")
}
