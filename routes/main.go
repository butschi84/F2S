package routes

import (
	"butschi84/f2s/logger"
	"log"
)

var logging *log.Logger

func init() {
	// initialize logging
	logging = logger.Initialize("routes")
}
