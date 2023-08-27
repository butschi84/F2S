package apiserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type F2SApiServerAuthBasic struct {
}

func (a *F2SApiServerAuthBasic) AuthenticateRequest(request *http.Request) error {
	// Get the Authorization header
	logging.Info(fmt.Sprintf("[auth] reading auth header"))

	// Check if Authorization header is present
	if request.Header.Get("Authorization") == "" {
		logging.Warn(fmt.Sprintf("[auth] authorization header missing"))
		return fmt.Errorf("Authorization header missing")
	}

	// Split the Authorization header into type and credentials
	authParts := strings.Split(request.Header.Get("Authorization"), " ")
	if len(authParts) != 2 || authParts[0] != "Basic" {
		logging.Warn(fmt.Sprintf("[auth] invalid authorization header format"))
		return fmt.Errorf("invalid authorization header format")
	}

	// Decode the base64 encoded credentials
	credentials, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		logging.Warn(fmt.Sprintf("[auth] error decoding credentials"))
		return fmt.Errorf("error decoding credentials")
	}

	// Split the credentials into username and password
	credentialsParts := strings.SplitN(string(credentials), ":", 2)
	if len(credentialsParts) != 2 {
		logging.Warn(fmt.Sprintf("[auth] invalid credentials format"))
		return fmt.Errorf("invalid credentials format")
	}

	username := credentialsParts[0]
	password := credentialsParts[1]

	logging.Info(fmt.Sprintf("[auth] find user '%s' credentials in config", username))
	for _, user := range f2shub.F2SConfiguration.Config.F2S.Auth.Basic {
		if user.Username == username {
			if user.Password == password {
				logging.Info(fmt.Sprintf("[auth] user '%s' authenticated successfully", username))

				// store information in request context
				ctx := context.WithValue(request.Context(), "username", username)
				ctx = context.WithValue(ctx, "usergroup", user.Group)
				*request = *request.WithContext(ctx)

				return nil
			} else {
				logging.Warn(fmt.Sprintf("[auth] password for user '%s' did not match", username))
				return fmt.Errorf("access denied")
			}
		}
	}

	logging.Warn(fmt.Sprintf("user '%s' not found in configuration. access denied", username))
	return fmt.Errorf("access denied")
}
