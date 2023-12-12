package eventmanager

import (
	"butschi84/f2s/services/logger"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var logging *logger.F2SLogger

func NewEventManager() *EventManager {

	// initialize logging
	logging = logger.Initialize("eventmanager")

	logging.Info("initializing new eventmanager")
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
	logging.Info(fmt.Sprintf("publish new event (type: '%s', uid: '%s', description: '%s')", string(event.Type), event.UID, event.Description))
	logging.Event(fmt.Sprintf("uid=%s type=%s description='%s'", event.UID, event.Type, event.Description))
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
	logging.Info("start eventhandling in async function")
	go func() {
		for {
			event := <-em.eventChannel
			for _, handler := range em.eventHandlers {
				go handler(event)
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

	logging.Info(fmt.Sprintf("add event '%s' to buffer", event.UID))

	// Shift existing events to the right
	n := len(e.LastEvents)
	logging.Info(fmt.Sprintf("current buffer length: %d", n))
	if n < 100 {
		logging.Debug("buffer not full yet. just append new event to buffer at position 0")
		// If the array is not full, increase its size
		e.LastEvents = append(e.LastEvents, PrettyEvent{})
		copy(e.LastEvents[1:], e.LastEvents[:n])
	} else {
		logging.Debug("buffer is full (100 items). drop oldest event and append new event at position 0")
		// If the array is already full, drop the last event
		copy(e.LastEvents[1:], e.LastEvents[:n-1])
	}
	// Add the new event at the beginning
	e.LastEvents[0] = prettyEvent
}
