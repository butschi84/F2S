package apiserver

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"butschi84/f2s/state/configuration"
	"fmt"
	"net/http"
	"regexp"

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
	protectedRouter.HandleFunc("/users/me", queryCurrentUser).Methods(http.MethodGet)
	protectedRouter.HandleFunc("/users", queryAllUsers).Methods(http.MethodGet)
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

		// user authentication
		authErr := authenticateIncomingRequest(request)
		if authErr != nil {
			logging.Warn(fmt.Sprintf("there was an error while authenticating the request: %s", authErr.Error()))
			http.Error(w, "access denied", http.StatusUnauthorized)
			return
		}

		// user authorization
		authType := f2shub.F2SConfiguration.Config.F2S.Auth.GlobalConfig.Type
		if authType != "" && authType != "none" {

			// parse subset of requestpath
			logging.Debug(fmt.Sprintf("parsing main request path from: %s", request.URL.Path))
			re := regexp.MustCompile(`/(?P<reqpath>[^/]*)`)
			matches := re.FindStringSubmatch(request.URL.Path)
			mainRequestPath := matches[re.SubexpIndex("reqpath")]
			logging.Debug(fmt.Sprintf("main req path: %s", mainRequestPath))

			// get request method
			requestMethod := request.Method

			// get users authorization group
			usergroup := request.Context().Value("usergroup").(string)
			username := request.Context().Value("username").(string)
			logging.Debug(fmt.Sprintf("user is in authorization group: %s", usergroup))

			var requiredPrivilege string
			if mainRequestPath == "invoke" {
				requiredPrivilege = string(configuration.F2SPrivilegeFunctionsInvoke)
			} else if mainRequestPath == "functions" && requestMethod == http.MethodGet {
				requiredPrivilege = string(configuration.F2SPrivilegeFunctionsList)
			} else if mainRequestPath == "functions" && requestMethod == http.MethodPost {
				requiredPrivilege = string(configuration.F2SPrivilegeFunctionsCreate)
			} else if mainRequestPath == "functions" && requestMethod == http.MethodDelete {
				requiredPrivilege = string(configuration.F2SPrivilegeFunctionsDelete)
			} else if mainRequestPath == "deployments" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "operator" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "services" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "endpoints" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "dispatcher" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "config" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "events" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "users" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "deployments" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			} else if mainRequestPath == "prometheus" {
				requiredPrivilege = string(configuration.F2SPrivilegeSettingsView)
			}

			// check if user has privilege
			if requiredPrivilege != "" && !configuration.HasGlobalPrivilege(requiredPrivilege, usergroup) {
				http.Error(w, fmt.Sprintf("access denied - missing role '%s' for user %s", requiredPrivilege, username), http.StatusUnauthorized)
				return
			}
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
