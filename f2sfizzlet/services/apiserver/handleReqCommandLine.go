package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

func handleReqCommandLine(w http.ResponseWriter, r *http.Request) {

	// get request method
	method := r.Method

	invocation := Invocation{
		File: "start.sh",
	}

	// Read the request body.
	if method == "POST" || method == "PUT" {
		logging.Info("parsing request body")

		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" {
			logging.Info("reuqest body is in json")
			var data map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			err :=
				decoder.Decode(&data)
			if err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			// Write the JSON data to a file named "input.json"
			jsonData, err := json.MarshalIndent(data, "", "    ")
			if err != nil {
				http.Error(w, "Error converting JSON data", http.StatusInternalServerError)
				return
			}

			// Generate a new random UUID
			newUUID := uuid.New()
			uuidString := newUUID.String()
			inputFilename := fmt.Sprintf("input-%s.json", uuidString)

			err = os.WriteFile(inputFilename, jsonData, 0644)
			if err != nil {
				http.Error(w, "Error writing to file", http.StatusInternalServerError)
				return
			}

			invocation.Parameters = inputFilename
			defer os.Remove(inputFilename)
		} else {
			// http.Error(w, "Content-Type is not JSON", http.StatusUnsupportedMediaType)

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}

			invocation.Parameters = string(body)
		}
	}

	// launch start script
	logging.Info("running start.sh")
	cmd := exec.Command("sh", invocation.File, invocation.Parameters)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		logging.Error(fmt.Errorf("Error executing the script: %s", err))
		http.Error(w, "Error Executing Command in Container", http.StatusInternalServerError)
		return
	}

	// set response headers

	var result map[string]interface{}
	if isJSON(stdout.Bytes()) {
		w.Header().Set("Content-Type", "application/json")
		err := json.Unmarshal([]byte(stdout.String()), &result)
		if err != nil {
			logging.Warn(fmt.Sprintf("failed to parse request result to json for request"))
		}
	} else {
		w.Header().Set("Content-Type", "text/plain")
		result = map[string]interface{}{"data": stdout.String()}
	}

	json.NewEncoder(w).Encode(RequestResult{
		Result: "ok",
		Data:   result,
	})
}

func isJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}
