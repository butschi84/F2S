package kafka

import (
	"butschi84/f2s/state/configuration"
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// start listening on a kafka topic
func initializeConsumer(kafkaListenerConfig *configuration.F2SConfigMapKafkaListener) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(f2shub.F2SConfiguration.Config.F2S.Kafka.Brokers, ","),
		"group.id":          kafkaListenerConfig.ConsumerGroup,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		logging.Info("Error creating consumer: %v\n", err.Error())
		time.Sleep(10 * time.Second) // Add a 10-second delay
		return err
	}

	defer consumer.Close()

	consumer.SubscribeTopics([]string{kafkaListenerConfig.Topic}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			logging.Info(fmt.Sprintf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value)))
			invokeActions := matchMessage(string(msg.Key), string(msg.Value), kafkaListenerConfig)
			for _, a := range invokeActions {
				for _, f2sFunctionUid := range a.F2SFunctions {
					logging.Info(fmt.Sprintf("adding invocation for function %s to queue", f2sFunctionUid))
					result, err := invokeFunction(f2sFunctionUid, string(msg.Value))

					// process result
					if err == nil {
						produceMessage(kafkaListenerConfig.Topic, "response-key", result)
					} else {
						produceMessage(kafkaListenerConfig.Topic, "response-key", err.Error())
					}
				}
			}

		} else {
			logging.Info(fmt.Sprintf("Error: %v\n", err))
		}
	}
}
