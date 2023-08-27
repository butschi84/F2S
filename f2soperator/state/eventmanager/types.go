package eventmanager

import "time"

// possible event types
type EventType string

const (
	Event_FunctionInvoked         EventType = "function invoked"
	Event_FunctionInvokationEnded EventType = "function invokation ended"
	Event_InflightRequestsChanged EventType = "inflight requests changed"
	Event_FunctionScaled          EventType = "function scaling changed"
	Event_ConfigurationChanged    EventType = "configuration changed"
	Event_EndpointsChanged        EventType = "endpoints changed"
)

type Event struct {
	// generate a random id (uid) for each event
	UID string
	// Data is the payload of the event
	Data interface{}
	// event type
	Type EventType
	// event description (optional)
	Description string
}

type PrettyEvent struct {
	UID         string    `json:"uid"`
	Type        EventType `json:"type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type EventHandler func(event Event)

type EventManager struct {
	eventChannel  chan Event
	eventHandlers []EventHandler
	LastEvents    []PrettyEvent // buffer that always contains 100 last events
}
