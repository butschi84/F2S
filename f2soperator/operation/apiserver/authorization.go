package apiserver

import (
	"butschi84/f2s/state/queue"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// get the current user
func getCurrentUser(request *http.Request) (queue.F2SAuthUser, error) {
	username, usernameOK := request.Context().Value("username").(string)
	groupname, groupnameOK := request.Context().Value("usergroup").(string)

	if !usernameOK {
		logging.Warn("there was an error while trying to read the current username")
		return queue.F2SAuthUser{}, fmt.Errorf("there was an error while trying to read the current username")
	}
	if !groupnameOK {
		logging.Warn("there was an error while trying to read the current users groupname")
		return queue.F2SAuthUser{}, fmt.Errorf("there was an error while trying to read the current users groupname")
	}

	return queue.F2SAuthUser{
		Username: username,
		Group:    groupname,
	}, nil
}

// get all users
func getUsers() ([]queue.F2SAuthUser, error) {

	logging.Info("request to query all users")
	var result []queue.F2SAuthUser

	switch f2shub.F2SConfiguration.Config.F2S.Auth.GlobalConfig.Type {
	case "basic":
		logging.Info("authentication is set to: basic")

		logging.Info(fmt.Sprintf("found %v users", len(f2shub.F2SConfiguration.Config.F2S.Auth.Basic)))
		result = make([]queue.F2SAuthUser, len(f2shub.F2SConfiguration.Config.F2S.Auth.Basic))

		for i, _ := range f2shub.F2SConfiguration.Config.F2S.Auth.Basic {
			result[i] = queue.F2SAuthUser{
				Username: f2shub.F2SConfiguration.Config.F2S.Auth.Basic[i].Username,
				Group:    f2shub.F2SConfiguration.Config.F2S.Auth.Basic[i].Group,
			}
		}
	case "token":
		logging.Info("authentication is set to: token")

		result = make([]queue.F2SAuthUser, len(f2shub.F2SConfiguration.Config.F2S.Auth.Token.Tokens))
		logging.Info(fmt.Sprintf("found %v users", len(f2shub.F2SConfiguration.Config.F2S.Auth.Token.Tokens)))

		jwtSecret := f2shub.F2SConfiguration.Config.F2S.Auth.Token.JwtSecret
		for i, _ := range f2shub.F2SConfiguration.Config.F2S.Auth.Token.Tokens {
			// decode token
			token, err := jwt.Parse(f2shub.F2SConfiguration.Config.F2S.Auth.Token.Tokens[i].Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil {
				logging.Error(fmt.Sprintf("error decoding token '%s': %s", f2shub.F2SConfiguration.Config.F2S.Auth.Token.Tokens[i].Token, err.Error()))
				continue
			}

			// decode claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				logging.Warn(fmt.Sprintf("invalid claims format in token '%s'", f2shub.F2SConfiguration.Config.F2S.Auth.Token.Tokens[i].Token))
				continue
			}

			// get username and group
			sub, subExists := claims["sub"].(string)
			group, groupExists := claims["group"].(string)

			if !subExists {
				logging.Warn(fmt.Sprintf("token '%s' has no 'sub' claim", f2shub.F2SConfiguration.Config.F2S.Auth.Token.Tokens[i].Token))
				continue
			}
			if !groupExists {
				logging.Warn(fmt.Sprintf("token '%s' has no 'group' claim", f2shub.F2SConfiguration.Config.F2S.Auth.Token.Tokens[i].Token))
				continue
			}

			result[i] = queue.F2SAuthUser{
				Username: sub,
				Group:    group,
			}
		}
	case "none":
		fallthrough
	default:
		result = make([]queue.F2SAuthUser, 0)
	}

	return result, nil
}
