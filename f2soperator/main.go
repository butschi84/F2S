package main

import (
	"butschi84/f2s/configuration"
	"butschi84/f2s/metrics"
	"butschi84/f2s/operator"
	"butschi84/f2s/routes"
	"butschi84/f2s/services/logger"
	"sync"
)

var F2SConfiguration configuration.F2SConfiguration
var logging logger.F2SLogger

func init() {
	// initialize logging
	logging = logger.Initialize("main")
}

func main() {

	logging.Info(" ")
	logging.Info("F2S Platform 0.0.1")
	logging.Info("=====================")

	logging.Info("=> initializing config package")
	F2SConfiguration = configuration.ActiveConfiguration

	var wg sync.WaitGroup

	// Number of goroutines to wait for
	numWorkers := 3
	wg.Add(numWorkers)

	// start api router
	logging.Info("=> initializng rest api server")
	go routes.HandleRequests(&F2SConfiguration, &wg)

	// start operator (manages deployments in f2s-containers namespace)
	logging.Info("=> initializng f2s-containers namespace operator")
	go operator.RunOperator(&F2SConfiguration, &wg)

	// start metrics
	logging.Info("=> initializng metrics")
	go metrics.HandleRequests(&F2SConfiguration, &wg)

	logging.Info("=> done initializing")
	wg.Wait()
}
