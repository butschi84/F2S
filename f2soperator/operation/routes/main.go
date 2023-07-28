package routes

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var logging logger.F2SLogger

var f2shub *hub.F2SHub

type Status struct {
	Status string `json:"status"`
}

func init() {
	// initialize logging
	logging = logger.Initialize("routes")
}

func HandleRequests(hub *hub.F2SHub, wg *sync.WaitGroup) {
	defer wg.Done()

	f2shub = hub
	router := mux.NewRouter().StrictSlash(false)

	router.Use(corsMiddleware)

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	f2shub.F2SEventManager.Subscribe(handleEvent)

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
	router.HandleFunc("/dispatcher", getCurrentDispatcherData).Methods(http.MethodGet)
	router.HandleFunc("/config", getConfiguration).Methods(http.MethodGet)
	router.HandleFunc("/events", getLatestEvents).Methods(http.MethodGet)
	router.HandleFunc("/deployments", createDeployment).Methods(http.MethodPost)
	router.HandleFunc("/invoke/{target}", invokeFunction)
	router.HandleFunc("/prometheus/{functionname}/{metricname}", getPrometheusMetric)
	router.HandleFunc("/health", checkHealth)

	// frontend, ui
	frontendHandler := http.FileServer(http.Dir("./static/frontend"))
	router.PathPrefix("/").Handler(frontendHandler)

	logging.Info("listening on http://0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", router)
}

// Middleware function to set Access-Control-Allow-Origin header to *
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Access-Control-Allow-Origin header to allow all origins (*)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests (OPTIONS).
		if r.Method == http.MethodOptions {
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
