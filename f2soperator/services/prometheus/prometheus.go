package prometheus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

// read current value of a prometheus metric
func ReadPrometheusMetric(metricName string, functionName string) (float64, error) {
	client := http.DefaultClient

	// Prepare the request URL with label selector
	requestURL := fmt.Sprintf("%s?query=%s{%s=\"%s\"}", "http://prometheus-service.f2s:9090/api/v1/query", "{__name__}", "functionname", functionName)
	// requestURL := fmt.Sprintf("%s?query=%s{%s=\"%s\"}", "http://192.168.2.40:32412/api/v1/query", metricName, "functionname", functionName)

	// Send GET request to Prometheus
	resp, err := client.Get(requestURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Parse the JSON response
	var promResponse PrometheusResponse
	if err := json.Unmarshal(body, &promResponse); err != nil {
		return 0, err
	}

	// Check if any metric result exists
	if len(promResponse.Data.Result) == 0 {
		return 0, fmt.Errorf("metric not found")
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
