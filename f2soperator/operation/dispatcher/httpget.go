package dispatcher

import (
	"butschi84/f2s/state/queue"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func httpGet(url string, result *queue.F2SRequestResult) error {
	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(f2shub.F2SConfiguration.Config.F2S.Timeouts.HttpTimeout),
	}

	// Send GET request using the client
	response, err := client.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("http_timeout: %s", err))
			return fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("error during httpGet function Call: %s", err))
			return err
		}

	}
	defer response.Body.Close()

	// Read response body
	fetchResponseError := fetchResponse(response, result)
	if fetchResponseError != nil {
		return fmt.Errorf("unable to fetch http response from container: %s", fetchResponseError.Error())
	}
	return nil
}

func httpPost(url, data string, result *queue.F2SRequestResult) error {
	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(f2shub.F2SConfiguration.Config.F2S.Timeouts.HttpTimeout),
	}

	// Create the request body as a byte buffer from the string data
	body := bytes.NewBufferString(data)

	// Send POST request using the client
	response, err := client.Post(url, "application/json", body)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("http_timeout: %s", err))
			return fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("error during httpPost function Call: %s", err))
			return err
		}
	}
	defer response.Body.Close()

	// Read response body
	fetchResponseError := fetchResponse(response, result)
	if fetchResponseError != nil {
		return fmt.Errorf("unable to fetch http response from container: %s", fetchResponseError.Error())
	}
	return nil
}

func httpPut(url, data string, result *queue.F2SRequestResult) error {
	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(f2shub.F2SConfiguration.Config.F2S.Timeouts.HttpTimeout),
	}

	// Create the request body as a byte buffer from the string data
	body := bytes.NewBufferString(data)

	// Create a PUT request using http.NewRequest
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		logging.Error(fmt.Errorf("error creating PUT request: %s", err))
		return err
	}
	// Set the appropriate content type for the request (e.g., "application/json" for JSON data)
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the client
	response, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("http_timeout: %s", err))
			return fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("error during httpPut function Call: %s", err))
			return err
		}
	}
	defer response.Body.Close()

	// Read response body
	fetchResponseError := fetchResponse(response, result)
	if fetchResponseError != nil {
		return fmt.Errorf("unable to fetch http response from container: %s", fetchResponseError.Error())
	}
	return nil
}

func httpDelete(url string, result *queue.F2SRequestResult) error {
	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(f2shub.F2SConfiguration.Config.F2S.Timeouts.HttpTimeout),
	}

	// Create a DELETE request using http.NewRequest
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		logging.Error(fmt.Errorf("error creating DELETE request: %s", err))
		return err
	}

	// Send the request using the client
	response, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("http_timeout: %s", err))
			return fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("error during httpDelete function Call: %s", err))
			return err
		}
	}
	defer response.Body.Close()

	// Read response body
	fetchResponseError := fetchResponse(response, result)
	if fetchResponseError != nil {
		return fmt.Errorf("unable to fetch http response from container: %s", fetchResponseError.Error())
	}
	return nil
}

// process response from container
func fetchResponse(response *http.Response, result *queue.F2SRequestResult) (err error) {
	// Read response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Error(fmt.Errorf("[fetchResponse] error when reading httpPost function call result body: %s", err))
		return err
	}

	result.Result = map[string]interface{}{"data": string(responseBody)}
	result.ContentType = "text/plain"

	// Check if the response is application/json
	if response.Header.Get("Content-Type") == "application/json" {
		logging.Debug("response is in json format")
		var parsedResult map[string]interface{}
		err := json.Unmarshal([]byte(responseBody), &parsedResult)
		if err != nil {
			logging.Warn(fmt.Sprintf("failed to parse request result to json for request: %s", &result.UID))
		}
		result.Result = parsedResult
		result.ContentType = response.Header.Get("Content-Type")
	}
	return nil
}
