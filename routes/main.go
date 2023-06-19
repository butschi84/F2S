package routes

import (
	config "butschi84/f2s/configuration"
	"butschi84/f2s/logger"
	"log"
	"net/http"

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

func HandleRequests(config config.F2SConfiguration) {
	F2SConfiguration = config
	router := mux.NewRouter().StrictSlash(true)

	// openAPI spec
	openAPIHandler := http.FileServer(http.Dir("./static/openapi"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", openAPIHandler))

	// retrieve configured f2s functions
	router.HandleFunc("/functions", returnAllFunctions)
	router.HandleFunc("/functions/{id}", getFunction)
	router.HandleFunc("/deployments", getAllDeployments)

	router.HandleFunc("/", root)

	http.ListenAndServe("localhost:8000", router)
}
