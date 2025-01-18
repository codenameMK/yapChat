package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"database/sql"
	_model "yap-chat/consumer/models/v1"
	postgresdb "yap-chat/postgres/v1"
	configenv "yap-chat/config/v1"

	"golang.org/x/term"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerStruct struct {
	Consumer *kafka.Consumer
	Topic    string
	pgdb	 *sql.DB
}

func (c *ConsumerStruct) StartConsumer(userId int32) {
	// Subscribe to the Kafka topic

	err := c.Consumer.Subscribe(c.Topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v\n", c.Topic, err)
	}

	log.Printf("Consuming messages from topic: %s\n", c.Topic)

	postgres := postgresdb.PostgresStruct{
		PostgresHost:     configenv.PostgresHost,
		PostgresPort:     configenv.PostgresPort,
		PostgresUser:     configenv.PostgresUser,
		PostgresPassword: configenv.PostgresPassword,
		PostgresDb:       configenv.PostgresDb,
	}

	pgdb, err := postgres.Init()

	for {
		var recMessage _model.MessageStruct
		// Read message from Kafka
		msg, err := c.Consumer.ReadMessage(-1)
		if err == nil {
			// Insert the message into PostgreSQL
			err := json.Unmarshal(msg.Value, &recMessage)
			if err != nil {
				log.Printf("Error unmarshalling message: %v\n", err)
				continue
			}
			if userId != recMessage.UserId && userId == recMessage.ReceiverId {
				err := postgresdb.InsertMessage(pgdb, recMessage.ReceiverId, time.Now(), recMessage.Message , recMessage.UserId, recMessage.ReceiverId)
				if err != nil {
					log.Fatalf("Error inserting message: %v", err)
				}
				
				printChatToRight(time.Now().Format("2006-01-02 15:04:05"))
				printChatToRight(recMessage.Message)
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

func printChatToRight(text string) error {
	// Get terminal width
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("error getting terminal size: %v", err)
	}

	// Calculate the right 60% of the screen
	rightStart := int(float64(width) * 0.4) // Start at 40% to use right 60%
	rightWidth := int(float64(width) * 0.6)

	// Split the text into words
	words := strings.Fields(text)
	currentLine := ""

	// Process each word
	for _, word := range words {
		// Check if adding this word would exceed the right width
		if len(currentLine)+len(word)+1 > rightWidth {
			// Print current line with proper padding
			padding := strings.Repeat(" ", rightStart)
			fmt.Printf("%s%s\n", padding, currentLine)
			currentLine = word
		} else {
			// Add word to current line
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	// Print last line if not empty
	if currentLine != "" {
		padding := strings.Repeat(" ", rightStart)
		fmt.Printf("%s%s\n", padding, currentLine)
	}

	return nil
}
