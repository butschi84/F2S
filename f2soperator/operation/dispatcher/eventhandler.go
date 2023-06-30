package dispatcher

import (
	"butschi84/f2s/state/eventmanager"
	"fmt"
)

func handleEvents(event eventmanager.Event) {
	logging.Info("processing event", fmt.Sprintf("'%s'", string(event.Type)))

}
