package main

import (
	"fmt"

	consumerHandler "yap-chat/consumer/handler/v1"
	producerHandler "yap-chat/producer/handler/v1"

	config "yap-chat/config/v1"
)

func main() {
	var userId int
	fmt.Print("Enter your user Id: ")
	fmt.Scanln(&userId)

	config.Init()

	topic := "comments"
	producerHandler := producerHandler.ProducerStruct{
		BrokerList:  config.KafkaBrokers,
		ClientId:    config.KafkaClientId,
		Topic:       topic,
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
	}
	// Conditionally add SSL fields only if SSL is enabled
	if config.KafkaSslEnabled {
		consumerHandler.SSLProtocol = "SSL"
		consumerHandler.SSLPath = config.KafkaSslCert
	}

	go consumerHandler.Init()

	select {}

}
