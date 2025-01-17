package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/term"

	configenv "yap-chat/config/v1"
	postgres "yap-chat/postgres/v1" 
)

type ConsumerStruct struct {
	BrokerList  string
	GroupId     int
	SSLProtocol string
	SSLPath     string
}

type MessageStruct struct {
	TimeStamp  time.Time
	Message    string
	UserId     int32
	ReceiverId int32
}

func (co *ConsumerStruct) Init() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": co.BrokerList,
		"group.id":          co.GroupId,
		"auto.offset.reset": "latest",
	}
	if co.SSLProtocol == "True" {
		config.SetKey("security.protocol", "SSL")
		config.SetKey("ssl.ca.location", co.SSLPath)
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer consumer.Close()

	// define postgres struct
	postgres := postgres.PostgresStruct{
		PostgresHost:     configenv.PostgresHost,
		PostgresPort:     configenv.PostgresPort,
		PostgresUser:     configenv.PostgresUser,
		PostgresPassword: configenv.PostgresPassword,
		PostgresDb:       configenv.PostgresDb,
	}

	pgdb, err := postgres.Init()

	StartConsumer(consumer, "comments", int32(co.GroupId), pgdb)
}

func StartConsumer(consumer *kafka.Consumer, topic string, userId int32, db *sql.DB) {
	// Subscribe to the Kafka topic

	err := consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v\n", topic, err)
	}

	log.Printf("Consuming messages from topic: %s\n", topic)

	for {
		var recMessage MessageStruct
		// Read message from Kafka
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			// Insert the message into PostgreSQL
			err := json.Unmarshal(msg.Value, &recMessage)
			if err != nil {
				log.Printf("Error unmarshalling message: %v\n", err)
				continue
			}
			if userId != recMessage.UserId {

				// Call InsertMessage, using fields directly from recmessage
				err := postgres.InsertMessage(db, recMessage.ReceiverId, time.Now(), recMessage.Message , recMessage.UserId, recMessage.ReceiverId)
				if err != nil {
					log.Fatalf("Error inserting message: %v", err)
				}
				
				PrintChatToRight(time.Now().Format("2006-01-02 15:04:05"))
				PrintChatToRight(recMessage.Message)
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

func PrintChatToRight(text string) error {
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
