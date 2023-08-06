package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func httpGet(url string) (string, error) {
	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(f2shub.F2SConfiguration.Config.F2S.Timeouts.HttpTimeout),
	}

	// Send GET request using the client
	response, err := client.Get(url)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("http_timeout: %s", err))
			return "", fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("error during httpGet function Call: %s", err))
			return "", err
		}

	}
	defer response.Body.Close()

	// Read response body
	resp, fetchResponseError := fetchResponse(response)
	if fetchResponseError != nil {
		return "", fmt.Errorf("unable to fetch http response from container: %s", fetchResponseError.Error())
	}
	return resp, nil
}

func httpPost(url, data string) (string, error) {
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
			return "", fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("error during httpPost function Call: %s", err))
			return "", err
		}
	}
	defer response.Body.Close()

	// Read response body
	resp, fetchResponseError := fetchResponse(response)
	if fetchResponseError != nil {
		return "", fmt.Errorf("unable to fetch http response from container: %s", fetchResponseError.Error())
	}
	return resp, nil
}

func httpPut(url, data string) (string, error) {
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
		return "", err
	}
	// Set the appropriate content type for the request (e.g., "application/json" for JSON data)
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the client
	response, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("http_timeout: %s", err))
			return "", fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("error during httpPut function Call: %s", err))
			return "", err
		}
	}
	defer response.Body.Close()

	// Read response body
	resp, fetchResponseError := fetchResponse(response)
	if fetchResponseError != nil {
		return "", fmt.Errorf("unable to fetch http response from container: %s", fetchResponseError.Error())
	}
	return resp, nil
}

func httpDelete(url string) (string, error) {
	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Millisecond * time.Duration(f2shub.F2SConfiguration.Config.F2S.Timeouts.HttpTimeout),
	}

	// Create a DELETE request using http.NewRequest
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		logging.Error(fmt.Errorf("error creating DELETE request: %s", err))
		return "", err
	}

	// Send the request using the client
	response, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			logging.Error(fmt.Errorf("http_timeout: %s", err))
			return "", fmt.Errorf("http_timeout: %s", err)
		} else {
			logging.Error(fmt.Errorf("error during httpDelete function Call: %s", err))
			return "", err
		}
	}
	defer response.Body.Close()

	// Read response body
	resp, fetchResponseError := fetchResponse(response)
	if fetchResponseError != nil {
		return "", fmt.Errorf("unable to fetch http response from container: %s", fetchResponseError.Error())
	}
	return resp, nil
}

// process response from container
func fetchResponse(response *http.Response) (resp string, err error) {
	// Read response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logging.Error(fmt.Errorf("[fetchResponse] error when reading httpPost function call result body: %s", err))
		return "", err
	}

	// Check if the response is application/json
	if response.Header.Get("Content-Type") == "application/json" {
		logging.Debug("response is in json format")
		parsedResponse, err := parseJSONRepsonse(string(responseBody))
		if err != nil {
			return "", fmt.Errorf("[fetchResponse] Could not parse json response body: %s", err.Error())
		}
		return parsedResponse, nil
	} else {
		// Convert response body to a string and return
		logging.Debug("response is not in json format")
		return string(responseBody), nil
	}
}

// convert response from example: `{\"data\":\"hello world\"}`
// to pretty json string:
//
//	{
//	  "data": "hello world"
//	}
func parseJSONRepsonse(response string) (parsedResponse string, err error) {
	unescapedResult, err := strconv.Unquote(`"` + response + `"`)
	if err != nil {
		return "", fmt.Errorf("[parseJSONRepsonse] Error unescaping JSON: %s. (response was: %s)", err.Error(), response)
	}

	// Unmarshal the JSON string into a map
	var dataMap map[string]interface{}
	err = json.Unmarshal([]byte(unescapedResult), &dataMap)
	if err != nil {
		return "", fmt.Errorf("[parseJSONRepsonse] Error encoding JSON: %s", err.Error())
	}

	// Marshal the dataMap into a pretty formatted JSON string
	prettyJSON, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("[parseJSONRepsonse] Error encoding JSON: %s", err.Error())
	}

	return string(prettyJSON), nil
}
