package v1

import (
	"database/sql"
	"fmt"
	"log"

	_HttpServer "yap-chat/consumer/v1"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	brokerList = "b-2.mskkafkacluster.gre9l6.c3.kafka.ap-south-1.amazonaws.com:9094,b-1.mskkafkacluster.gre9l6.c3.kafka.ap-south-1.amazonaws.com:9094" // Use port 9094 for TLS
	topic      = "comments"
	caCertPath = "producer/resources/AmazonRootCA1.pem" // Replace with the actual path to the CA certificate
	groupID    = "comment-consumer-group"               // Consumer group ID for Kafka
	insertQuery = "INSERT INTO kafka_messages (message) VALUES ($1)" // Query to insert messages
	postgreshost = "yapchat-db.c30quw6ambjb.ap-south-1.rds.amazonaws.com"
	postgresport = 5432
)

// Init initializes the Kafka consumer with SSL configuration, connects to PostgreSQL, and starts consuming messages
func Init() {

	// Build the PostgreSQL connection string
	dbConnStr := fmt.Sprintf("postgreshost=%s postgresport=%s user=yourusername password=yourpassword dbname=yourdbname sslmode=disable", postgreshost, postgresport)

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v\n", err)
	}
	defer db.Close()

	// Kafka consumer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokerList,
		"group.id":          groupID,
		"security.protocol": "SSL",
		"ssl.ca.location":   caCertPath,
		"auto.offset.reset": "latest", // Start consuming from the earliest message if no offsets are present
	}

	// Create a new Kafka consumer
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v\n", err)
	}
	defer consumer.Close()

	// Initialize KafkaConsumer with PostgreSQL connection and insert query
	kafkaConsumer := _HttpServer.NewKafkaConsumer(consumer, topic, db, insertQuery)

	// Start consuming messages
	kafkaConsumer.StartConsumer()
}