package v1

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	KafkaBrokers          string
	KafkaSslEnabled       bool
	KafkaSecurityProtocol string
	KafkaSslCert          string
	KafkaTopic            string
	KafkaClientId         string
)

func Init() {
	// Load environment variables from .env (optional)
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	KafkaBrokers = os.Getenv("KAFKA_BROKERS")
	KafkaSslEnabled, err = strconv.ParseBool(os.Getenv("KAFKA_SSL_ENABLED"))
	if err != nil {
		log.Fatalf("Error converting string to bool: %v", err)
	}
	KafkaSecurityProtocol = os.Getenv("KAFKA_SECURITY_PROTOCOL")
	KafkaSslCert = os.Getenv("KAFKA_SSL_CERT")
	KafkaTopic = os.Getenv("KAFKA_TOPIC")
	KafkaClientId = os.Getenv("KAFKA_CLIENT_ID")
}
