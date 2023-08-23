package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func authenticateIncomingRequest(request *http.Request) error {
	switch f2shub.F2SConfiguration.Config.F2S.Auth.GlobalConfig.Type {
	case "basic":
		auth := F2SApiServerAuthBasic{}
		return auth.AuthenticateRequest(request)
	case "token":
		auth := F2SApiServerAuthToken{}
		return auth.AuthenticateRequest(request)
	case "none":
		fallthrough
	default:
		auth := F2SApiServerAuthNone{}
		return auth.AuthenticateRequest(request)
	}
}

type JWTIssuerRequestBody struct {
	JWTSecret string `json:"jwtsecret"`
	Username  string `json:"username"`
	Group     string `json:"group"`
}

// *********************************************************
// all endpoints
// *********************************************************
func signJWT(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to sign JWT")

	// Parse JSON request body
	var requestBody JWTIssuerRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		return
	}

	// Access the attributes
	jwtSecret := requestBody.JWTSecret
	username := requestBody.Username
	group := requestBody.Group

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   username,
		"group": group,
		"exp":   time.Now().Add(time.Hour * 24 * 365).Unix(), // Token expires in 1 year
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "Failed to sign JWT", http.StatusInternalServerError)
		return
	}

	// Create and send the JSON response
	jsonResponse := map[string]string{"token": tokenString}
	responseJSON, err := json.Marshal(jsonResponse)
	if err != nil {
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(responseJSON)
}

func getAuthType(w http.ResponseWriter, r *http.Request) {
	logging.Info("request to get auth type")
	response := f2shub.F2SConfiguration.Config.F2S.Auth.GlobalConfig.Type
	logging.Info(fmt.Sprintf("configured auth type is: %s", response))

	// Set the content type to plain text
	w.Header().Set("Content-Type", "text/plain")

	// Write the string response to the response writer
	_, _ = w.Write([]byte(response))
}

// this is a protected route. so just send 'ok' response
func signin(w http.ResponseWriter, r *http.Request) {
	response := SimpleHealthResponse{Status: "ok"}

	json.NewEncoder(w).Encode(response)
}
