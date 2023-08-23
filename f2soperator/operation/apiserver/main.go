package apiserver

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var logging logger.F2SLogger

var f2shub *hub.F2SHub

type Status struct {
	Status string `json:"status"`
}

func init() {
	// initialize logging
	logging = logger.Initialize("apiserver")
}

func HandleRequests(hub *hub.F2SHub) {

	f2shub = hub
	router := mux.NewRouter().StrictSlash(false)

	router.Use(corsMiddleware)

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	f2shub.F2SEventManager.Subscribe(handleEvent)

	// openAPI spec
	openAPIHandler := http.FileServer(http.Dir("./static/openapi"))
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", openAPIHandler))

	// open endpoints
	router.HandleFunc("/health", checkHealth)
	router.HandleFunc("/auth/type", getAuthType).Methods(http.MethodGet)
	router.HandleFunc("/auth/signjwt", signJWT)

	// protected endpoints (auth)
	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(authMiddleware)
	protectedRouter.HandleFunc("/auth/signin", signin)
	protectedRouter.HandleFunc("/functions", returnAllFunctions).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/functions", createFunction).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/functions/{id}", getFunction).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/functions/{id}", deleteFunction).Methods(http.MethodDelete)
	protectedRouter.HandleFunc("/deployments", getAllDeployments).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/operator", getOperatorState).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/services", getAllServices).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/endpoints", getAllEndpoints).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/dispatcher", getCurrentDispatcherData).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/config", getConfiguration).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/events", getLatestEvents).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/deployments", createDeployment).Methods(http.MethodPost)
	protectedRouter.HandleFunc("/invoke/{target}", invokeFunction)
	protectedRouter.HandleFunc("/prometheus/{functionname}/{metricname}", getPrometheusMetric)

	// frontend, ui
	frontendHandler := http.FileServer(http.Dir("./static/frontend"))
	router.PathPrefix("/").Handler(frontendHandler)

	logging.Info("listening on http://0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", router)
}

// Authentication Middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {

		logging.Info(fmt.Sprintf("authenticate request with '%s' auth model", f2shub.F2SConfiguration.Config.F2S.Auth.GlobalConfig.Type))
		authErr := authenticateIncomingRequest(request)
		if authErr != nil {
			logging.Warn(fmt.Sprintf("there was an error while authenticating the request: %s", authErr.Error()))
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, request)
	})
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
