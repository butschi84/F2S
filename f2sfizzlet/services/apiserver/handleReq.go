package apiserver

import (
	"net/http"
)

// *********************************************************
// all functions
// *********************************************************
func handleRequest(w http.ResponseWriter, r *http.Request) {
	logging.Info("running request")

	mode := isPortListening()

	switch mode {
	case ModeReverseProxy:
		logging.Info("running in reverse proxy mode")
		handleRequestRverseProxy(w, r)
	case ModeCommandLine:
		logging.Info("running in command line mode")
		handleReqCommandLine(w, r)
	}
}
