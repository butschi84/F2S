package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// *********************************************************
// simple health endpoint
// *********************************************************
func queryAllUsers(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get all users")

	response, err := getUsers()
	if err != nil {
		logging.Warn(fmt.Sprintf("failed to query users: %s", err.Error()))
	}

	json.NewEncoder(w).Encode(response)
}

func queryCurrentUser(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get current user")

	response, err := getCurrentUser(r)
	if err != nil {
		logging.Warn(fmt.Sprintf("failed to query current user: %s", err.Error()))
	}

	json.NewEncoder(w).Encode(response)
}
