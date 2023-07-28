package eventmanager

import (
	"time"

	"github.com/google/uuid"
)

func NewEventManager() *EventManager {
	eventmanager := &EventManager{
		eventChannel: make(chan Event),
		LastEvents:   make([]PrettyEvent, 0, 100),
	}
	eventmanager.Start()

	// add internal handler for buffer / cache
	eventmanager.Subscribe(eventmanager.bufferAddEvent)

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

func (e *EventManager) bufferAddEvent(event Event) {
	prettyEvent := PrettyEvent{
		UID:         event.UID,
		Type:        event.Type,
		Description: event.Description,
		Timestamp:   time.Now(),
	}

	// Shift existing events to the right
	n := len(e.LastEvents)
	if n < 100 {
		// If the array is not full, increase its size
		e.LastEvents = append(e.LastEvents, PrettyEvent{})
		copy(e.LastEvents[1:], e.LastEvents[:n])
	} else {
		// If the array is already full, drop the last event
		copy(e.LastEvents[1:], e.LastEvents[:n-1])
	}
	// Add the new event at the beginning
	e.LastEvents[0] = prettyEvent
}
