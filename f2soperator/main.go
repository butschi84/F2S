package main

import (
	"butschi84/f2s/configuration"
	"butschi84/f2s/dispatcher"
	"butschi84/f2s/hub"
	"butschi84/f2s/metrics"
	"butschi84/f2s/operator"
	"butschi84/f2s/routes"
	"butschi84/f2s/services/logger"
	"sync"
)

var F2SConfiguration configuration.F2SConfiguration
var logging logger.F2SLogger

var F2SHub hub.F2SHub

func init() {
	// initialize logging
	logging = logger.Initialize("main")
}

func main() {

	logging.Info(" ")
	logging.Info("F2S Platform 0.0.1")
	logging.Info("=====================")

	logging.Info("=> initializing config package")
	F2SHub = hub.F2SHub{
		F2SConfiguration: configuration.ActiveConfiguration,
	}

	var wg sync.WaitGroup

	// Number of goroutines to wait for
	numWorkers := 4
	wg.Add(numWorkers)

	// start api router
	logging.Info("=> initializng rest api server")
	go routes.HandleRequests(&F2SHub, &wg)

	// start operator (manages deployments in f2s-containers namespace)
	logging.Info("=> initializng f2s-containers namespace operator")
	go operator.RunOperator(&F2SHub, &wg)

	// start metrics
	logging.Info("=> initializng metrics")
	go metrics.HandleRequests(&F2SHub, &wg)

	// start dispatcher
	dispatcher := dispatcher.F2SDispatcher{
		Hub:       &F2SHub,
		WaitGroup: &wg,
	}
	go dispatcher.HandleRequests()

	logging.Info("=> done initializing")
	wg.Wait()
}
