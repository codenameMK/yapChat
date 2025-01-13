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
		SSLProtocol: func() string { if config.KafkaSslEnabled { return "SSL" } else { return "" } }(),
		SSLPath:     func() string { if config.KafkaSslEnabled { return config.KafkaSslCert } else { return "" } }(),
		Topic:       topic,
	}
	go producerHandler.Init()

	//consumer
	consumerHandler := consumerHandler.ConsumerStruct{
		BrokerList:  config.KafkaBrokers,
		GroupId:     userId,
		SSLProtocol: func() string { if config.KafkaSslEnabled { return "SSL" } else { return "" } }(),
		SSLPath:     func() string { if config.KafkaSslEnabled { return config.KafkaSslCert } else { return "" } }(),
	}
	go consumerHandler.Init()

	select {}

}
