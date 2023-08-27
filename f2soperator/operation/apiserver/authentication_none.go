package apiserver

import (
	"context"
	"net/http"
)

type F2SApiServerAuthNone struct {
}

func (a *F2SApiServerAuthNone) AuthenticateRequest(request *http.Request) error {

	// store information in request context
	ctx := context.WithValue(request.Context(), "username", "anonymous")
	ctx = context.WithValue(ctx, "usergroup", "anonymous")
	*request = *request.WithContext(ctx)

	return nil
}
