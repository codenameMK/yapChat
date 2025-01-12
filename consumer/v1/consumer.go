package v1

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

type KafkaConsumer struct {
	Consumer      *kafka.Consumer
	Topic         string
	DbConnection  *sql.DB // PostgreSQL connection
	InsertQuery   string  // SQL query to insert the message
}

func (kc *KafkaConsumer) StartConsumer() {
	// Subscribe to the Kafka topic
	err := kc.Consumer.Subscribe(kc.Topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v\n", kc.Topic, err)
	}

	// Start a Goroutine to consume messages
	go func() {
		log.Printf("Consuming messages from topic: %s\n", kc.Topic)

		for {
			// Read message from Kafka
			msg, err := kc.Consumer.ReadMessage(-1)
			if err == nil {
				// Insert the message into PostgreSQL
				err = kc.insertMessage(msg.Value)
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

	}()

}

func (kc *KafkaConsumer) insertMessage(message []byte) error {
	// Prepare the SQL query to insert the message into PostgreSQL
	// Assume we have a table `kafka_messages` with columns `message` (text) and `created_at` (timestamp)
	// _, err := kc.DbConnection.Exec(kc.InsertQuery, string(message))
	// if err != nil {
	// 	return fmt.Errorf("could not insert message: %v", err)
	// }
	fmt.Println("Message received :> %s\n", string(message))

	// Optionally log the success of the insert
	log.Printf("Message inserted into PostgreSQL: %s\n", string(message))
	return nil
}

func NewKafkaConsumer(consumer *kafka.Consumer, topic string, dbConn *sql.DB, insertQuery string) *KafkaConsumer {
	return &KafkaConsumer{
		Consumer:     consumer,
		Topic:        topic,
		DbConnection: dbConn,
		InsertQuery:  insertQuery,
	}
}