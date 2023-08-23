package apiserver

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type F2SApiServerAuthToken struct {
}

func (a *F2SApiServerAuthToken) AuthenticateRequest(request *http.Request) error {
	// Get the token from the request headers, assuming it's in the "Authorization" header
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		logging.Warn(fmt.Sprintf("[auth] authorization header missing"))
		return fmt.Errorf("Authorization header missing")
	}

	// Extract the token string (assuming it's in the format "Bearer <token>")
	tokenString := authHeader[len("Bearer "):]

	// Parse the token using the provided secret key
	jwtSecret := f2shub.F2SConfiguration.Config.F2S.Auth.Token.JwtSecret
	logging.Warn(fmt.Sprintf("[auth] jwt secret: %s", jwtSecret))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		logging.Warn(fmt.Sprintf("[auth] failed to parse JWT"))
		return fmt.Errorf("failed to parse JWT")
	}

	// Check if the token is valid
	if !token.Valid {
		logging.Warn(fmt.Sprintf("[auth] invalid jwt token"))
		return fmt.Errorf("invalid jwt token")
	}

	// Access the claims (including the "sub" attribute)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logging.Warn(fmt.Sprintf("[auth] invalid claims format"))
		return fmt.Errorf("invalid claims format")
	}

	// Access the "sub" attribute
	sub, subExists := claims["sub"].(string)
	if !subExists {
		logging.Warn(fmt.Sprintf("[auth] sub claim not found"))
		return fmt.Errorf("sub claim not found")
	}

	// Use the sub attribute
	logging.Info(fmt.Sprintf("[auth] sub: %s", sub))

	return nil
}
