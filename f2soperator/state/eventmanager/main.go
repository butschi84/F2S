package eventmanager

import (
	typesV1alpha1 "butschi84/f2s/state/configuration/api/types/v1alpha1"

	"github.com/google/uuid"
)

// possible event types
type EventType string

const (
	Event_FunctionInvoked         EventType = "function invoked"
	Event_FunctionInvokationEnded EventType = "function invokation ended"
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
	// container for a f2sfunction object
	Function typesV1alpha1.Function
}

type EventHandler func(event Event)

type EventManager struct {
	eventChannel  chan Event
	eventHandlers []EventHandler
}

func NewEventManager() *EventManager {
	eventmanager := &EventManager{
		eventChannel: make(chan Event),
	}
	eventmanager.Start()
	return eventmanager
}

// function to publish a new event on eventmanager
func (em *EventManager) Publish(event Event) {
	em.eventChannel <- event
}

// function to generate a random uuid
func (em *EventManager) GenerateUUID() string {
	// Generate a new random UUID
	uuid := uuid.New()

	return uuid.String()
}

// function to subscribe to events from eventmanager
func (em *EventManager) Subscribe(handler EventHandler) {
	em.eventHandlers = append(em.eventHandlers, handler)
}

// start event manager
func (em *EventManager) Start() {
	go func() {
		for {
			event := <-em.eventChannel
			for _, handler := range em.eventHandlers {
				handler(event)
			}
		}
	}()
}
