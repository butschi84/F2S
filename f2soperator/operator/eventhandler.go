package operator

import (
	eventmanager "butschi84/f2s/services/eventmanager"
	"fmt"
)

func handleEvent(event eventmanager.Event) {
	logging.Info("processing event", fmt.Sprintf("'%s'", string(event.Type)))

	switch event.Type {
	case eventmanager.Event_ConfigurationChanged:
		if master {
			logging.Info("configuration has changed. rebalance immediately")
			Rebalance()
		}
	}
}
