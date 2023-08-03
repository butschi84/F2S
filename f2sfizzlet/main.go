package main

import (
	"butschi84/f2sfizzlet/services/apiserver"
	"butschi84/f2sfizzlet/services/f2slogger"
	"sync"
)

func main() {

	logging := f2slogger.Initialize("main")

	logging.Info(" ")
	logging.Info("F2S Fizzlet 0.0.1")
	logging.Info("=====================")

	logging.Info("=> initializing")
	var wg sync.WaitGroup

	// Number of goroutines to wait for
	numWorkers := 1
	wg.Add(numWorkers)

	// start api router
	logging.Info("=> initializng rest api server")
	go apiserver.HandleRequests(&wg)

	logging.Info("=> done initializing")
	wg.Wait()
}
