package v1

import (
	"log"
	server "yap-chat/consumer/v1"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerStruct struct {
	BrokerList  string
	GroupId     int
	SSLProtocol string
	SSLPath     string
	Topic       string
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
	startConsumer := server.ConsumerStruct{
		Consumer: consumer,
		Topic:    co.Topic,
	}
	startConsumer.StartConsumer(int32(co.GroupId))
}
