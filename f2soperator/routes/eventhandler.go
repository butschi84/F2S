package routes

import (
	"butschi84/f2s/services/eventmanager"
	"fmt"
)

func handleEvent(event eventmanager.Event) {
	logging.Info("processing event", fmt.Sprintf("'%s'", string(event.Type)))
}
