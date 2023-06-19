package routes

import "butschi84/f2s/services/eventmanager"

func handleEvent(event eventmanager.Event) {
	logging.Println("processing event", event)
}
