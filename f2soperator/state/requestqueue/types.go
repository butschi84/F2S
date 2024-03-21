package queue

import (
	"butschi84/f2s/state/configuration/api/types/v1alpha1"
	"time"
)

// possible event types
type EventType string

const (
	Event_IncomingRequest     EventType = "incoming request"
	Event_RequestCompleted    EventType = "incoming completed"
	Event_MinimumAvailability EventType = "function scaled for minimum availability"
)

type F2SAuthUser struct {
	Username string
	Group    string
}

// F2SRequest to invoke a function
type F2SInvocationRequest struct {
	Metadata F2SInvocationRequestMetadata
	Request  F2SInvocationRequestRequest
	Target   F2SInvocationRequestTarget
	Response F2SInvocationRequestResponse
}

type F2SInvocationRequestMetadata struct {
	Uid      string
	StartAt  time.Time
	EndAt    time.Time
	Duration float32
}
type F2SInvocationRequestRequest struct {
	Path     string
	Method   string
	Payload  string
	Function v1alpha1.PrettyFunction
	F2SUser  F2SAuthUser
}
type F2SInvocationRequestTarget struct {
	Uid              string
	Address          string
	InflightRequests int
	Name             string
	Url              string
}

type F2SInvocationRequestResponse struct {
	Success     bool                   `json:"success"`
	Result      map[string]interface{} `json:"result"`
	Message     string                 `json:"details"`
	ContentType string                 `json:"contentType"`
}

type RequestHandler func(request *F2SInvocationRequest)
type F2SQueue struct {
	Requests []F2SInvocationRequest

	// dispatcher will subscribe to new requests
	eventChannel  chan *F2SInvocationRequest
	eventHandlers []RequestHandler
}
type IF2SQueue interface {
	AddRequest(req F2SInvocationRequest)
}
