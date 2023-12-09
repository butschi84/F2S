package kafka

import (
	"butschi84/f2s/state/configuration"
	"fmt"
	"strings"
)

// function to check if a configured kafka trigger matches the message
// => input (key, value)
// compare to f2s config
func matchMessage(key string, value string, kafkaListenerConfig *configuration.F2SConfigMapKafkaListener) []configuration.F2SConfigMapKafkaListenerAction {
	logging.Info(fmt.Sprintf("Kafka Message: %s => %s", key, value))

	// prepare results
	resInvokeActions := make([]configuration.F2SConfigMapKafkaListenerAction, 0)

	for _, action := range kafkaListenerConfig.Actions {
		for _, trigger := range action.Triggers {

			// get the actual message value or key that should be matched by this trigger
			var actualValue string = ""
			switch trigger.Type {
			case "key":
				actualValue = key
			case "value":
				actualValue = value
			}

			// see if it matches
			switch trigger.Filter {
			case "equal":
				if actualValue == trigger.Value {
					resInvokeActions = append(resInvokeActions, action)
				}
			case "not equal":
				if actualValue != trigger.Value {
					resInvokeActions = append(resInvokeActions, action)
				}
			case "contains":
				if strings.Contains(actualValue, trigger.Value) {
					resInvokeActions = append(resInvokeActions, action)
				}
			case "not contains":
				if !strings.Contains(actualValue, trigger.Value) {
					resInvokeActions = append(resInvokeActions, action)
				}
			}
		}
	}

	logging.Info(fmt.Sprintf("%v triggers matched the message.", len(resInvokeActions)))
	return resInvokeActions
}
