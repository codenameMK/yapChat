package v1

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	_model "yap-chat/producer/models/v1"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const maxMessageSize = 256

type ProducerStruct struct {
	Producer *kafka.Producer
	Topic    *string
}

func Validate(message string) error {
	if len(message) > maxMessageSize {
		return errors.New("message exceeds maximum allowed size")
	}
	return nil
}

func (p *ProducerStruct) StartServer(userId int32, receiverId int32) {

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
		err = Validate(text)
		if err != nil {
			log.Printf("Message is too large: %v\n", err)
			continue
		}
		message := _model.MessageStruct{
			TimeStamp:  time.Now(),
			Message:    text,
			UserId:     userId,     // Replace with actual user ID
			ReceiverId: receiverId, // Replace with actual sender ID
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message to JSON: %v\n", err)
			continue
		}

		// Produce the message to the Kafka topic
		err = p.Producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: p.Topic, Partition: kafka.PartitionAny},
			Value:          messageBytes,
		}, nil)

		if err != nil {
			fmt.Println("Failed to produce message: %v\n", err)
		} else {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		}
	}

}
