package operator

import (
	"butschi84/f2s/configuration"
	"butschi84/f2s/logger"
	"sync"
	"time"

	"log"
)

var logging *log.Logger

func init() {
	// initialize logging
	logging = logger.Initialize("operator")

}

func RunOperator(config *configuration.F2SConfiguration, wg *sync.WaitGroup) {
	defer wg.Done()

	// subscribe to configuration changes
	logging.Println("subscribing to config package events")
	config.EventManager.Subscribe(handleEvent)

	for {
		// Perform the desired task
		logging.Println("rebalancing...")

		// Sleep for 30 seconds
		time.Sleep(10 * time.Second)
	}
}

// manage k8s deployments in namespace f2s-containers
func Rebalance() {
	logging.Println("starting rebalance")

	// check for surplus deployments in f2s-containers namespace
	// functions := configuration.ActiveConfiguration.Functions
	// deployments :=
}
