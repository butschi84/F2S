package queue

import (
	"butschi84/f2s/services/logger"
	"fmt"
)

var logging logger.F2SLogger

func init() {
	// initialize logging
	logging = logger.Initialize("queue")
}

func Initialize() *F2SQueue {
	f2squeue := F2SQueue{
		Requests:     make([]F2SRequest, 0),
		eventChannel: make(chan F2SRequest),
	}
	logging.Info("starting queue")
	f2squeue.Start()

	return &f2squeue
}

func (f2squeue *F2SQueue) Start() {
	// start channel
	go func() {
		for {
			event := <-f2squeue.eventChannel
			logging.Info("processing new event")
			for _, handler := range f2squeue.eventHandlers {
				go handler(event)
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
	logging.Info("adding request to queue")
	queue.eventChannel <- req
}

// add a new request to the queue
func (queue *F2SQueue) RequestDone(req F2SRequest) {
	logging.Info(fmt.Sprintf("request %s has completed", req.UID))
	for x, request := range queue.Requests {
		if request.UID == req.UID {
			queue.Requests = append(queue.Requests[:x], queue.Requests[x+1:]...)
			return
		}
	}
}

// empty the queue
func (queue *F2SQueue) Clear(req F2SQueue) {

}
