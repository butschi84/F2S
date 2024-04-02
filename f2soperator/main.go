package main

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/operation/apiserver"
	"butschi84/f2s/operation/cluster"
	"butschi84/f2s/operation/dispatcher"
	"butschi84/f2s/operation/kafka"
	"butschi84/f2s/operation/metrics"
	"butschi84/f2s/operation/operator"
	"butschi84/f2s/services/logger"
	clusterstate "butschi84/f2s/state/cluster"
	"butschi84/f2s/state/configuration"
	"butschi84/f2s/state/dispatcherstate"
	"butschi84/f2s/state/eventmanager"
	"butschi84/f2s/state/operatorstate"
	"butschi84/f2s/state/queue"
	"fmt"

	"sync"
)

var (
	F2SConfiguration configuration.F2SConfiguration
	logging          *logger.F2SLogger
	F2SHub           hub.F2SHub
)

func init() {
	// initialize logging
	logging = logger.Initialize("main")
}

func handleComponent(name string, f func(*hub.F2SHub), wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		logging.Info(fmt.Sprintf("%s: Starting", name))
		f(&F2SHub)
		logging.Info(fmt.Sprintf("%s: Exited. Restarting...", name))
	}
}

func main() {

	logging.Info(" ")
	logging.Info("F2S Platform 0.0.1")
	logging.Info("=====================")

	logging.Info("=> initializing f2shub")
	F2SHub = hub.F2SHub{
		F2SEventManager:  eventmanager.NewEventManager(),
		F2SConfiguration: configuration.Initialize(),
		F2SQueue:         queue.Initialize(),
		F2SDispatcherHub: dispatcherstate.Initialize(),
		F2SOperatorState: operatorstate.Initialize(),
		F2SClusterState:  clusterstate.Initialize(),
	}

	var wg sync.WaitGroup

	// Number of goroutines to wait for
	numWorkers := 4
	wg.Add(numWorkers)

	// start all components
	logging.Info("=> initializng f2s cluster")
	go handleComponent("cluster", cluster.Initialize, &wg)
	logging.Info("=> initializng rest api server")
	go handleComponent("api server", apiserver.HandleRequests, &wg)
	go logging.Info("=> initializing operator")
	go handleComponent("operator", operator.RunOperator, &wg)
	logging.Info("=> initializng metrics")
	go handleComponent("metrics", metrics.HandleRequests, &wg)
	logging.Info("=> initializng request dispatcher")
	go handleComponent("dispatcher", dispatcher.Initialize, &wg)
	logging.Info("=> initializng kafka integration")
	go handleComponent("kafka", kafka.Initialize, &wg)

	logging.Info("=> done initializing")
	wg.Wait()
}
