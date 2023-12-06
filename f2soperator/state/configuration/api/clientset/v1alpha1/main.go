package v1alpha1

import (
	"butschi84/f2s/services/logger"

	"golang.org/x/exp/slog"
)

var logging *slog.Logger

func init() {
	// initialize logging
	logging = logger.Initialize("v1alpha1")
}
