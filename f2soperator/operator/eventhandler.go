package operator

import (
	eventmanager "butschi84/f2s/services/eventmanager"
)

func handleEvent(event eventmanager.Event) {
	logging.Println("processing event", event)

	switch event.Type {
	case eventmanager.Event_ConfigurationChanged:
		if master {
			logging.Println("configuration has changed. rebalance immediately")
			Rebalance()
		}
	}
}
