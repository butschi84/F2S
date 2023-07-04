package queue

// F2SRequest to invoke a function
type F2SRequest struct {
	UID           string
	Path          string
	Method        string
	ResultChannel chan F2SRequestResult
}

type F2SRequestResult struct {
	UID     string
	Success bool
	Result  string
}

type RequestHandler func(request F2SRequest)
type F2SQueue struct {
	Requests []F2SRequest

	// dispatcher will subscribe to new requests
	eventChannel  chan F2SRequest
	eventHandlers []RequestHandler
}
type IF2SQueue interface {
	AddRequest(req F2SRequest)
}
