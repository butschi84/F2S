package kubernetesservice

import (
	"butschi84/f2s/services/logger"

	"golang.org/x/exp/slog"
)

var logging *slog.Logger

func init() {
	// initialize logging
	logging = logger.Initialize("kubernetesservice")
}

func int32Ptr(i int32) *int32 {
	return &i
}
