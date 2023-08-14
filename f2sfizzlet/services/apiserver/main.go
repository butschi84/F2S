package apiserver

import (
	"butschi84/f2sfizzlet/services/f2slogger"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var logging f2slogger.F2SLogger

type Status struct {
	Status string `json:"status"`
}

func init() {
	// initialize logging
	logging = f2slogger.Initialize("apiserver")
}

func HandleRequests(wg *sync.WaitGroup) {
	defer wg.Done()

	router := mux.NewRouter().StrictSlash(false)

	router.Use(corsMiddleware)

	router.HandleFunc("/", handleRequest)

	// frontend, ui
	frontendHandler := http.FileServer(http.Dir("./static/frontend"))
	router.PathPrefix("/").Handler(frontendHandler)

	logging.Info("listening on http://0.0.0.0:9092")
	http.ListenAndServe("0.0.0.0:9092", router)
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
