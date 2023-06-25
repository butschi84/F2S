package prometheus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func ReadPrometheusMetricValue(metricName string, labels map[string]string) (float64, error) {

	promResponse, err := ReadPrometheusMetric(metricName, labels)
	if err != nil {
		return 0, err
	}

	// Extract the metric value
	value := promResponse.Data.Result[0].Value[1]

	// Parse the metric value as a float64
	metricValue, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	if err != nil {
		return 0, err
	}

	return metricValue, nil
}

// read current value of a prometheus metric
func ReadPrometheusMetric(metricName string, labels map[string]string) (PrometheusResponse, error) {
	client := http.DefaultClient

	// Prepare the label selector
	var labelSelectors []string
	for key, value := range labels {
		labelSelectors = append(labelSelectors, fmt.Sprintf("%s=\"%s\"", key, value))
	}
	labelSelector := strings.Join(labelSelectors, ",")

	// Prepare the request URL with label selector
	requestURL := fmt.Sprintf("http://prometheus-service.f2s:9090/api/v1/query?query=%s{%s}", metricName, labelSelector)
	// requestURL := fmt.Sprintf("http://192.168.2.40:32412/api/v1/query?query=%s{%s}", metricName, labelSelector)
	// requestURL := fmt.Sprintf("%s?query=%s{%s=\"%s\"}", "http://192.168.2.40:32412/api/v1/query", metricName, "functionname", functionName)

	// Send GET request to Prometheus
	resp, err := client.Get(requestURL)
	if err != nil {
		return PrometheusResponse{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PrometheusResponse{}, err
	}

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