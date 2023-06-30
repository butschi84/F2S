package queue

var Queue F2SQueue

func init() {
	f2squeue := F2SQueue{
		Requests:     make([]F2SRequest, 0),
		eventChannel: make(chan F2SRequest),
	}

	// start channel
	go func() {
		for {
			event := <-f2squeue.eventChannel
			for _, handler := range f2squeue.eventHandlers {
				handler(event)
			}
		}
	}()
}

// function to subscribe to events from eventmanager
func (em *F2SQueue) Subscribe(handler RequestHandler) {
	em.eventHandlers = append(em.eventHandlers, handler)
}

// add a new request to the queue
func (queue *F2SQueue) AddRequest(req F2SRequest) {
	Queue.eventChannel <- req
}

// empty the queue
func (queue *F2SQueue) Clear(req F2SQueue) {

}
