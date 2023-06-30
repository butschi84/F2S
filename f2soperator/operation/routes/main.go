package routes

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var logging logger.F2SLogger

var F2SHub hub.F2SHub

type Status struct {
	Status string `json:"status"`
}

func init() {
	// initialize logging
	logging = logger.Initialize("routes")
}

func HandleRequests(hub *hub.F2SHub, wg *sync.WaitGroup) {
	defer wg.Done()

	F2SHub = *hub
	router := mux.NewRouter().StrictSlash(false)

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	F2SHub.F2SEventManager.Subscribe(handleEvent)

	// openAPI spec
	openAPIHandler := http.FileServer(http.Dir("./static/openapi"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", openAPIHandler))

	// retrieve configured f2s functions
	router.HandleFunc("/functions", returnAllFunctions).Methods(http.MethodGet)
	router.HandleFunc("/functions", createFunction).Methods(http.MethodPost)
	router.HandleFunc("/functions/{id}", getFunction).Methods(http.MethodGet)
	router.HandleFunc("/functions/{id}", deleteFunction).Methods(http.MethodDelete)
	router.HandleFunc("/deployments", getAllDeployments).Methods(http.MethodGet)
	router.HandleFunc("/services", getAllServices).Methods(http.MethodGet)
	router.HandleFunc("/endpoints", getAllEndpoints).Methods(http.MethodGet)
	router.HandleFunc("/config", getConfiguration).Methods(http.MethodGet)
	router.HandleFunc("/deployments", createDeployment).Methods(http.MethodPost)
	router.HandleFunc("/invoke/{target}", invokeFunction)
	router.HandleFunc("/prometheus/{functionname}/{metricname}", getPrometheusMetric)

	router.HandleFunc("/", root)
	logging.Info("listening on http://0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", router)
}
