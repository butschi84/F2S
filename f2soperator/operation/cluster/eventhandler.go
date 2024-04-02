package cluster

import (
	"butschi84/f2s/state/eventmanager"
	"fmt"
)

// handle eventmanager events
func handleEvents(event eventmanager.Event) {
	logging.Info("processing event", fmt.Sprintf("'%s'", string(event.Type)))

	switch event.Type {
	}
}
