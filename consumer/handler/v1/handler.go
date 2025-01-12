package v1

import (
	"log"
	consumerv1 "yap-chat/consumer/v1"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)
type KafkaConsumer struct {
	UserID int
}

const (
	brokerList = "b-2.mskkafkacluster.gre9l6.c3.kafka.ap-south-1.amazonaws.com:9094,b-1.mskkafkacluster.gre9l6.c3.kafka.ap-south-1.amazonaws.com:9094" // Use port 9094 for TLS
	topic      = "comments"
	caCertPath = "consumer/resources/AmazonRootCA1.pem" // Replace with the actual path to the CA certificate
)

func (c *KafkaConsumer)Init(){
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokerList,
		"group.id":          c.UserID,
		"auto.offset.reset": "latest",
		"security.protocol": "SSL",
		"ssl.ca.location":   caCertPath,
		"debug":             "security,broker,consumer",
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()
	server := consumerv1.ConsumerHandler{
		Consumer: consumer,
		Topic:    topic,
	}
	server.StartConsuming()

}