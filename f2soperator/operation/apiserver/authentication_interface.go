package apiserver

import "butschi84/f2s/state/queue"

type IF2SApiServerAuth interface {
	AuthenticateRequest(request *queue.F2SRequest) error
}
