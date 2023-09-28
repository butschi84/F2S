package kafka

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"fmt"
)

var logging logger.F2SLogger
var f2shub *hub.F2SHub

func init() {
	// initialize logging
	logging = logger.Initialize("kafka")
}

func Initialize(hub *hub.F2SHub) {
	f2shub = hub

	// check if kafka integration is enabled in config
	if !hub.F2SConfiguration.Config.F2S.Kafka.Enabled {
		logging.Info("kafka integration is disabled via config")
		return
	}

	// get all configured kafka listeners
	logging.Info("query all configured kafka listeners from config")
	kafkaListeners := f2shub.F2SConfiguration.Config.F2S.Kafka.Listeners
	logging.Info(fmt.Sprintf("%v listeners are configured", len(kafkaListeners)))

	for _, consumer := range kafkaListeners {
		initializeConsumer(&consumer)
	}

}
