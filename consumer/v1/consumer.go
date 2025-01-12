package v1

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)


type ConsumerHandler struct {
	Consumer *kafka.Consumer
	Topic    string
}


func (ch *ConsumerHandler) StartConsuming() {
	// Subscribe to the topic
	err := ch.Consumer.Subscribe(ch.Topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v", ch.Topic, err)
	}

	fmt.Println("Consumer started, waiting for messages...")

	// Consume messages in a loop
	for {
		msg, err := ch.Consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Print the consumed message
		fmt.Printf("Message received: %s\n", string(msg.Value))
	}
}