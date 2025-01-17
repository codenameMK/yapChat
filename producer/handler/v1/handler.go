package v1

import (
	"log"
	server "yap-chat/producer/v1"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ProducerStruct struct {
	BrokerList  string
	ClientId    string
	SSLProtocol string
	SSLPath     string
	Topic       string
	UserId      int32
	ReceiverId  int32
}

func (po *ProducerStruct) Init() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": po.BrokerList,
		"client.id":         po.ClientId,
	}
	if po.SSLProtocol == "True" {
		config.SetKey("security.protocol", "SSL")
		config.SetKey("ssl.ca.location", po.SSLPath)
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()
	startProducer := server.ProducerStruct{
		Producer: producer,
		Topic:    &po.Topic,
	}
	startProducer.StartServer(po.UserId, po.ReceiverId)
}
