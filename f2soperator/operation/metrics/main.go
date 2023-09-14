package metrics

import (
	"butschi84/f2s/hub"
	"butschi84/f2s/services/logger"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logging logger.F2SLogger

// pointer to F2SConfiguration
var f2shub hub.F2SHub

// metrics
var metricTotalIncomingRequests *prometheus.CounterVec
var metricTotalCompletedRequests *prometheus.CounterVec
var metricTotalRequestDuration *prometheus.CounterVec
var metricTotalFailedRequests *prometheus.CounterVec
var metricActiveRequests *prometheus.GaugeVec
var metricLastRequestCompletion *prometheus.GaugeVec
var metricRequestDuration *prometheus.HistogramVec
var metricLastFunctionScaling *prometheus.GaugeVec

func init() {
	// initialize logging
	logging = logger.Initialize("metrics")

	// metric - total incoming requests
	metricTotalIncomingRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "f2s_requests_incoming_total",
			Help: "Total number of incoming requests",
		},
		[]string{"target", "functionuid", "functionname", "username"},
	)

	// metric - total request duration
	metricTotalRequestDuration = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "f2s_requests_duration_total",
			Help: "Total duration of all completed requests",
		},
		[]string{"target", "functionuid", "functionname", "username", "podUID"},
	)

	// metric - total completed requests
	metricTotalCompletedRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "f2s_requests_completed_total",
			Help: "Total number of completed requests",
		},
		[]string{"target", "functionuid", "functionname", "podUID"},
	)

	// metric - total failed requests
	metricTotalFailedRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "f2s_requests_failed_total",
			Help: "Total number of failed requests",
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

	// metric - lsat function scaling
	metricLastFunctionScaling = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "f2s_function_last_scaling_timestamp",
			Help: "Timestamp of last function scaling",
		},
		[]string{"target", "functionuid", "functionname"},
	)

	// metric - request duration
	// buckets specification for measuring response time
	// from 10ms to 60s ¯\_(ツ)_/¯
	// 10ms, 50ms, 100ms, 250ms, 500ms, 1s, 2s, 5s, 10s, 30s, 1m, 2m, 5m
	buckets := []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0, 10.0, 30.0, 60.0, 120.0, 300.0}
	metricRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "f2s_request_duration_seconds",
			Help:    "Histogram of request duration",
			Buckets: buckets,
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

	prometheus.MustRegister(metricTotalIncomingRequests)
	prometheus.MustRegister(metricTotalRequestDuration)
	prometheus.MustRegister(metricTotalCompletedRequests)
	prometheus.MustRegister(metricTotalFailedRequests)
	prometheus.MustRegister(metricActiveRequests)
	prometheus.MustRegister(metricLastRequestCompletion)
	prometheus.MustRegister(metricRequestDuration)
	prometheus.MustRegister(metricLastFunctionScaling)
}

func HandleRequests(hub *hub.F2SHub) {

	f2shub = *hub

	// subscribe to configuration changes
	logging.Info("subscribing to config package events")
	hub.F2SEventManager.Subscribe(handleEvent)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe("0.0.0.0:8081", router)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	handler := promhttp.Handler()
	handler.ServeHTTP(w, r)
}
