package eventmanager

type Event struct {
	// Data is the payload of the event
	Data interface{}
}

type EventHandler func(event Event)

type EventManager struct {
	eventChannel  chan Event
	eventHandlers []EventHandler
}

func NewEventManager() *EventManager {
	return &EventManager{
		eventChannel: make(chan Event),
	}
}

func (em *EventManager) Publish(event Event) {
	em.eventChannel <- event
}

func (em *EventManager) Subscribe(handler EventHandler) {
	em.eventHandlers = append(em.eventHandlers, handler)
}

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
