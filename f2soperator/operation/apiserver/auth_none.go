package apiserver

import "net/http"

type F2SApiServerAuthNone struct {
}

func (a *F2SApiServerAuthNone) AuthenticateRequest(request *http.Request) error {
	return nil
}
