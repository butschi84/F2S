package dispatcher

import (
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
