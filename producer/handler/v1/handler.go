package v1

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerStruct struct {
	BrokerList  string
	ClientId    string
	SSLProtocol string
	SSLPath     string
	Topic       string
}

func (po *ProducerStruct) Init() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": po.BrokerList,
		"client.id":         po.ClientId,
		"security.protocol": po.SSLProtocol,
		"ssl.ca.location":   po.SSLPath,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()
	StartServer(producer, &po.Topic)
}

func StartServer(producer *kafka.Producer, topic *string) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter messages to publish to the Kafka topic (type 'exit' to quit):")

	for {
		// Read input from terminal
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v\n", err)
			continue
		}

		// Trim whitespace and check for exit command
		text = text[:len(text)-1] // Remove the newline character
		if text == "exit" {
			break
		}

		// Produce the message to the Kafka topic
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
			Value:          []byte(text),
		}, nil)

		if err != nil {
			fmt.Println("Failed to produce message: %v\n", err)
		} else {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		}
	}

}
