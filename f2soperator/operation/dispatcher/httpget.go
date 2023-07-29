package dispatcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.Error(fmt.Errorf("error when reading httpGet functioncall result body: %s", err))
		return "", err
	}

	// Print response body
	return string(body), nil
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
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.Error(fmt.Errorf("error when reading httpPost function call result body: %s", err))
		return "", err
	}

	// Convert response body to a string and return
	return string(responseBody), nil
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
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.Error(fmt.Errorf("error when reading httpPut function call result body: %s", err))
		return "", err
	}

	// Convert response body to a string and return
	return string(responseBody), nil
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
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.Error(fmt.Errorf("error when reading httpDelete function call result body: %s", err))
		return "", err
	}

	// Convert response body to a string and return
	return string(responseBody), nil
}
