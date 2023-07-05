package queue

// F2SRequest to invoke a function
type F2SRequest struct {
	UID           string
	Path          string
	Method        string
	ResultChannel chan F2SRequestResult
}

type F2SRequestResult struct {
	UID     string `json:"uid"`
	Success bool   `json:"success"`
	Result  string `json:"result"`

	Request F2SRequest `json:"-"`

	Duration                   float64 `json:"duration"`
	DurationPerInflightRequest float64 `json:"duration_per_inflightrequest"`
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
