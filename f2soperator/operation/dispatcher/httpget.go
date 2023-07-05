package dispatcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpGet(url string) (string, error) {
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	defer response.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Print response body
	return string(body), nil
}
