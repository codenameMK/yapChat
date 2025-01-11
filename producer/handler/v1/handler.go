package v1

import (
	"log"
	_HttpServer "yap-chat/producer/v1"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	brokerList = "b-2.mskkafkacluster.gre9l6.c3.kafka.ap-south-1.amazonaws.com:9094,b-1.mskkafkacluster.gre9l6.c3.kafka.ap-south-1.amazonaws.com:9094" // Use port 9094 for TLS
	topic      = "comments"
	caCertPath = "producer/resources/AmazonRootCA1.pem" // Replace with the actual path to the CA certificate
)

func Init(){
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokerList,
		"client.id":         "go-producer",
		"security.protocol": "SSL",
		"ssl.ca.location":   caCertPath,
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()
	server := _HttpServer.HttpServer{
		Producer: producer,
		Topic:    "comments",
	}
	server.StartServer()

}