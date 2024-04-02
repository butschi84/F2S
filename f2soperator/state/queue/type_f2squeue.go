package queue

type RequestHandler func(request *F2SRequest)
type F2SQueue struct {
	Requests []F2SRequest

	// dispatcher will subscribe to new requests
	eventChannel  chan *F2SRequest
	eventHandlers []RequestHandler
}
type IF2SQueue interface {
	AddRequest(req F2SRequest)
}
