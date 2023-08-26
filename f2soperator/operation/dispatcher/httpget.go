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
			logging.Error(fmt.Errorf("[%s] http_timeout: %s", result.Request.UID, err))
			return fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("[%s] error during httpGet function Call: %s", result.Request.UID, err))
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

func isJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}

func httpPost(url, data string, result *queue.F2SRequestResult) error {
	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(f2shub.F2SConfiguration.Config.F2S.Timeouts.HttpTimeout),
	}

	// Create the request body as a byte buffer from the string data
	body := bytes.NewBufferString(data)

	// determine content type
	contentType := "text/plain"
	if isJSON(body.Bytes()) {
		contentType = "application/json"
	}
	logging.Info(fmt.Sprintf("[%s] determined post content-type: %s", result.Request.UID, contentType))

	// Send POST request using the client
	response, err := client.Post(url, contentType, body)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("http_timeout: %s", err))
			return fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("[%s] error during httpPost function Call: %s", result.Request.UID, err))
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
		logging.Error(fmt.Errorf("[%s] error creating PUT request: %s", result.Request.UID, err))
		return err
	}

	// determine content type
	contentType := "text/plain"
	if isJSON(body.Bytes()) {
		contentType = "application/json"
	}
	logging.Info(fmt.Sprintf("[%s] determined put content-type: %s", result.Request.UID, contentType))
	req.Header.Set("Content-Type", contentType)

	// Send the request using the client
	response, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("[%s] http_timeout: %s", result.Request.UID, err))
			return fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("[%s] error during httpPut function Call: %s", result.Request.UID, err))
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
		logging.Error(fmt.Errorf("[%s] error creating DELETE request: %s", result.Request.UID, err))
		return err
	}

	// Send the request using the client
	response, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("[%s] http_timeout: %s", result.Request.UID, err))
			return fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("[%s] error during httpDelete function Call: %s", result.Request.UID, err))
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
		logging.Error(fmt.Errorf("[%s] [fetchResponse] error when reading httpPost function call result body: %s", result.Request.UID, err))
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
			logging.Warn(fmt.Sprintf("[%s] failed to parse request result to json for request: %s", result.Request.UID, &result.UID))
		}
		result.Result = parsedResult
		result.ContentType = response.Header.Get("Content-Type")
	}
	return nil
}
