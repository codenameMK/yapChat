package main

import (
	"fmt"

	consumerHandler "yap-chat/consumer/handler/v1"
	producerHandler "yap-chat/producer/handler/v1"

	config "yap-chat/config/v1"
)

func main() {
	var userId, receiverId int
	fmt.Print("Enter your user Id: ")
	fmt.Scanln(&userId)
	fmt.Print("Enter your sender Id: ")
	fmt.Scanln(&receiverId)

	config.Init()

	topic := "comments"
	producerHandler := producerHandler.ProducerStruct{
		BrokerList: config.KafkaBrokers,
		ClientId:   config.KafkaClientId,
		Topic:      topic,
		UserId:     int32(userId),
		ReceiverId: int32(receiverId),
	}
	// Conditionally add SSL fields only if SSL is enabled
	if config.KafkaSslEnabled {
		producerHandler.SSLProtocol = "SSL"
		producerHandler.SSLPath = config.KafkaSslCert
	}
	go producerHandler.Init()

	// Initialize consumerHandler without SSL fields
	consumerHandler := consumerHandler.ConsumerStruct{
		BrokerList: config.KafkaBrokers,
		GroupId:    userId,
		Topic:      topic,
	}
	// Conditionally add SSL fields only if SSL is enabled
	if config.KafkaSslEnabled {
		consumerHandler.SSLProtocol = "SSL"
		consumerHandler.SSLPath = config.KafkaSslCert
	}

	go consumerHandler.Init()

	select {}

}
