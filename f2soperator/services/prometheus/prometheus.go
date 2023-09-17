package prometheus

import (
	"butschi84/f2s/state/configuration"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// get the most recent metric value
func ReadPrometheusMetricValue(config *configuration.F2SConfiguration, metricQuery string) (float64, error) {

	promResponse, err := ReadPrometheusMetric(config, metricQuery)
	if err != nil {
		return 0, err
	}

	if len(promResponse.Data.Result) == 0 {
		return 0.0, fmt.Errorf("could not read metric value. metric not found")
	}

	// Extract the metric value
	value := promResponse.Data.Result[0].Values[len(promResponse.Data.Result[0].Values)-1][1]

	// Parse the metric value as a float64
	metricValue, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	if err != nil {
		return 0, err
	}

	return metricValue, nil
}

// read current value of a prometheus metric
func ReadPrometheusMetric(config *configuration.F2SConfiguration, queryString string) (PrometheusResponse, error) {
	client := http.DefaultClient

	// Get the current time (now)
	now := time.Now()

	// Calculate the time two hours ago
	twoHoursAgo := now.Add(-2 * time.Hour)

	// Prepare the request URL
	encodedQueryString := url.QueryEscape(queryString)
	requestURL := fmt.Sprintf("http://%s/api/v1/query_range?query=%s&start=%d&end=%d&step=1m", config.Config.Prometheus.URL, encodedQueryString, twoHoursAgo.Unix(), now.Unix())

	// Send GET request to Prometheus
	resp, err := client.Get(requestURL)
	if err != nil {
		return PrometheusResponse{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PrometheusResponse{}, err
	}

	// fmt.Println(string(body))

	// Parse the JSON response
	var promResponse PrometheusResponse
	if err := json.Unmarshal(body, &promResponse); err != nil {
		return PrometheusResponse{}, err
	}

	// Check if any metric result exists
	if len(promResponse.Data.Result) == 0 {
		return PrometheusResponse{}, fmt.Errorf("metric not found")
	}

	return promResponse, nil
}
