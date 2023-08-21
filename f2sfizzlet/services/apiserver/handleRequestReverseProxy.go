package apiserver

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func handleRequestRverseProxy(w http.ResponseWriter, r *http.Request) {
	// Create a new URL for the target (where you want to forward the request)
	targetURL, err := url.Parse("http://127.0.0.1:8080")
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Update the request URL with the target host
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	// Perform the reverse proxy
	proxy.ServeHTTP(w, r)
}
