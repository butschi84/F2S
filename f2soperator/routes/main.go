package routes

import (
	config "butschi84/f2s/configuration"
	"butschi84/f2s/logger"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var logging *log.Logger

// pointer to F2SConfiguration
var F2SConfiguration config.F2SConfiguration

type Status struct {
	Status string `json:"status"`
}

func init() {
	// initialize logging
	logging = logger.Initialize("routes")
}

func HandleRequests(config *config.F2SConfiguration, wg *sync.WaitGroup) {
	defer wg.Done()

	F2SConfiguration = *config
	router := mux.NewRouter().StrictSlash(false)

	// subscribe to configuration changes
	logging.Println("subscribing to config package events")
	config.EventManager.Subscribe(handleEvent)

	// openAPI spec
	openAPIHandler := http.FileServer(http.Dir("./static/openapi"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", openAPIHandler))

	// retrieve configured f2s functions
	router.HandleFunc("/functions", returnAllFunctions).Methods(http.MethodGet)
	router.HandleFunc("/functions", createFunction).Methods(http.MethodPost)
	router.HandleFunc("/functions/{id}", getFunction).Methods(http.MethodGet)
	router.HandleFunc("/functions/{id}", deleteFunction).Methods(http.MethodDelete)
	router.HandleFunc("/deployments", getAllDeployments).Methods(http.MethodGet)
	router.HandleFunc("/deployments", createDeployment).Methods(http.MethodPost)
	router.HandleFunc("/invoke/{target}", invokeFunction)

	router.HandleFunc("/", root)

	http.ListenAndServe("0.0.0.0:8000", router)
}
