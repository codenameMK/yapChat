package v1

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type HttpServer struct {
	Producer *kafka.Producer
	Topic    string
}

func (s HttpServer) StartServer() {
	// Channel to capture errors and graceful shutdown signals
	done := make(chan bool)

	// Start a Goroutine to read messages from the terminal and publish to Kafka
	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter messages to publish to the Kafka topic (type 'exit' to quit):")

		for {
			// Read input from terminal
			fmt.Print("> ")
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
			err = s.Producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &s.Topic, Partition: kafka.PartitionAny},
				Value:          []byte(text),
			}, nil)

			if err != nil {
				log.Printf("Failed to produce message: %v\n", err)
			} else {
				log.Printf("Message '%s' queued for delivery\n", text)
			}
		}

		// Signal to stop the server
		done <- true
	}()

	// Wait for completion or shutdown signal
	<-done

	// Ensure all messages are delivered before shutting down
	log.Println("Flushing pending messages...")
	s.Producer.Flush(15 * 1000)

	log.Println("Server shutting down...")
}
