package main

import (
	"butschi84/f2s/configuration"
	"butschi84/f2s/logger"
	"butschi84/f2s/operator"
	"butschi84/f2s/routes"
	"log"
	"sync"
)

var F2SConfiguration configuration.F2SConfiguration
var logging *log.Logger

func init() {
	// initialize logging
	logging = logger.Initialize("main")
}

func main() {

	logging.Println(" ")
	logging.Println("F2S Platform 0.0.1")
	logging.Println("=====================")

	logging.Println("=> initializing config package")
	F2SConfiguration = configuration.ActiveConfiguration

	var wg sync.WaitGroup

	// Number of goroutines to wait for
	numWorkers := 2
	wg.Add(numWorkers)

	// start api router
	logging.Println("=> initializng rest api server")
	go routes.HandleRequests(&F2SConfiguration, &wg)

	// start operator (manages deployments in f2s-containers namespace)
	logging.Println("=> initializng f2s-containers namespace operator")
	go operator.RunOperator(&F2SConfiguration, &wg)

	logging.Println("=> done initializing")
	wg.Wait()
}
