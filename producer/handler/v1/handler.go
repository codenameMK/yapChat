package v1

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const maxMessageSize = 256

type ProducerStruct struct {
	BrokerList  string
	ClientId    string
	SSLProtocol string
	SSLPath     string
	Topic       string
	UserId      int32
}

func (m *MessageStruct) Validate() error {
	if len(m.Message) > maxMessageSize {
		return errors.New("message exceeds maximum allowed size")
	}
	return nil
}

type MessageStruct struct {
	TimeStamp time.Time
	Message   string
	UserId    int32
	SenderId  int32
}

func (po *ProducerStruct) Init() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": po.BrokerList,
		"client.id":         po.ClientId,
	}
	if po.SSLProtocol == "True" {
		config.SetKey("security.protocol", "SSL")
		config.SetKey("ssl.ca.location", po.SSLPath)
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()
	StartServer(producer, &po.Topic, po.UserId)
}

func StartServer(producer *kafka.Producer, topic *string, userId int32) {

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
		message := MessageStruct{
			TimeStamp: time.Now(),
			Message:   text,
			UserId:    userId, // Replace with actual user ID
			SenderId:  456,    // Replace with actual sender ID
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message to JSON: %v\n", err)
			continue
		}

		// Produce the message to the Kafka topic
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
			Value:          messageBytes,
		}, nil)

		if err != nil {
			fmt.Println("Failed to produce message: %v\n", err)
		} else {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		}
	}

}
