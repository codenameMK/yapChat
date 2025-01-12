package v1

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerStruct struct {
	BrokerList  string
	GroupId     int
	SSLProtocol string
	SSLPath     string
}

func (co *ConsumerStruct) Init() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": co.BrokerList,
		"group.id":          co.GroupId,
		"auto.offset.reset": "latest",
		"security.protocol": "SSL",
		"ssl.ca.location":   co.SSLPath,
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer consumer.Close()
	StartConsumer(consumer, "comments")
}

func StartConsumer(consumer *kafka.Consumer, topic string) {
	// Subscribe to the Kafka topic
	err := consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v\n", topic, err)
	}

	log.Printf("Consuming messages from topic: %s\n", topic)

	for {
		// Read message from Kafka
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			// Insert the message into PostgreSQL
			fmt.Println(string(msg.Value))
			if err != nil {
				log.Printf("Error inserting message into PostgreSQL: %v\n", err)
			}
		} else {
			// Check if it's a non-fatal error
			kafkaErr, isKafkaError := err.(kafka.Error)
			if isKafkaError && kafkaErr.IsFatal() {
				log.Printf("Fatal Kafka error: %v\n", kafkaErr)
				break
			}
			log.Printf("Error reading message: %v\n", err)
		}
	}
}
