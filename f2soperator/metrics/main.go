package metrics

import (
	config "butschi84/f2s/configuration"
	"butschi84/f2s/logger"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logging *log.Logger

// pointer to F2SConfiguration
var F2SConfiguration config.F2SConfiguration

// metrics
var metricTotalRequests *prometheus.CounterVec
var metricActiveRequests *prometheus.GaugeVec
var metricLastRequestCompletion *prometheus.GaugeVec
var metricRequestDuration *prometheus.HistogramVec

func init() {
	// initialize logging
	logging = logger.Initialize("metrics")

	// metric - total requests
	metricTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "f2s_requests_total",
			Help: "Total number of requests",
		},
		[]string{"target", "functionuid", "functionname"},
	)

	// metric - active requests
	metricActiveRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "f2s_active_requests_total",
			Help: "Total number of currently active requests",
		},
		[]string{"target", "functionuid", "functionname"},
	)

	// metric - request duration
	metricRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "f2s_request_duration_seconds",
			Help:    "Histogram of request duration",
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
		},
		[]string{"target", "functionuid", "functionname"},
	)

	// metric - timestamp of last request completion
	metricLastRequestCompletion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "f2s_last_request_timestamp",
			Help: "Timestamp of the last request",
		},
		[]string{"target", "functionuid", "functionname"},
	)

	prometheus.MustRegister(metricTotalRequests)
	prometheus.MustRegister(metricActiveRequests)
	prometheus.MustRegister(metricLastRequestCompletion)
	prometheus.MustRegister(metricRequestDuration)
}

func HandleRequests(config *config.F2SConfiguration, wg *sync.WaitGroup) {
	defer wg.Done()

	F2SConfiguration = *config

	// subscribe to configuration changes
	logging.Println("subscribing to config package events")
	config.EventManager.Subscribe(handleEvent)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe("0.0.0.0:8081", router)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	handler := promhttp.Handler()
	handler.ServeHTTP(w, r)
}
