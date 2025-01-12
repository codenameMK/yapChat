package main

import (
	"fmt"

	consumerHandler "yap-chat/consumer/handler/v1"
	producerHandler "yap-chat/producer/handler/v1"
)

const (
	brokerList = "b-2.mskkafkacluster.gre9l6.c3.kafka.ap-south-1.amazonaws.com:9094,b-1.mskkafkacluster.gre9l6.c3.kafka.ap-south-1.amazonaws.com:9094" // Use port 9094 for TLS
	topic      = "comments"
	caCertPath = "resources/AmazonRootCA1.pem" // Replace with the actual path to the CA certificate
	protocol   = "SSL"
	clientId   = "producer"
)

// func StartConsumer(consumer *kafka.Consumer, topic string) {
// 	// Subscribe to the Kafka topic
// 	err := consumer.Subscribe(topic, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to subscribe to topic %s: %v\n", topic, err)
// 	}

// 	log.Printf("Consuming messages from topic: %s\n", topic)

// 	for {
// 		// Read message from Kafka
// 		msg, err := consumer.ReadMessage(-1)
// 		if err == nil {
// 			// Insert the message into PostgreSQL
// 			fmt.Println(string(msg.Value))
// 			if err != nil {
// 				log.Printf("Error inserting message into PostgreSQL: %v\n", err)
// 			}
// 		} else {
// 			// Check if it's a non-fatal error
// 			kafkaErr, isKafkaError := err.(kafka.Error)
// 			if isKafkaError && kafkaErr.IsFatal() {
// 				log.Printf("Fatal Kafka error: %v\n", kafkaErr)
// 				break
// 			}
// 			log.Printf("Error reading message: %v\n", err)
// 		}
// 	}
// }

// func StartServer(producer *kafka.Producer, topic *string) {

// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Println("Enter messages to publish to the Kafka topic (type 'exit' to quit):")

// 	for {
// 		// Read input from terminal
// 		fmt.Print("> ")
// 		text, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Printf("Error reading input: %v\n", err)
// 			continue
// 		}

// 		// Trim whitespace and check for exit command
// 		text = text[:len(text)-1] // Remove the newline character
// 		if text == "exit" {
// 			break
// 		}

// 		// Produce the message to the Kafka topic
// 		err = producer.Produce(&kafka.Message{
// 			TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
// 			Value:          []byte(text),
// 		}, nil)

// 		if err != nil {
// 			log.Printf("Failed to produce message: %v\n", err)
// 		} else {
// 			log.Printf("Message '%s' queued for delivery\n", text)
// 		}
// 	}

// }

func main() {
	var userId int
	fmt.Print("Enter your user Id: ")
	fmt.Scanln(&userId)

	topic := "comments"
	producerHandler := producerHandler.ProducerStruct{
		BrokerList:  brokerList,
		ClientId:    clientId,
		SSLProtocol: protocol,
		SSLPath:     caCertPath,
		Topic:       topic,
	}
	go producerHandler.Init()

	//consumer
	consumerHandler := consumerHandler.ConsumerStruct{
		BrokerList:  brokerList,
		GroupId:     userId,
		SSLProtocol: protocol,
		SSLPath:     caCertPath,
	}
	go consumerHandler.Init()

	select {}
	// consumerHandlernow := consuerHandlerv1.KafkaConsumer{UserID: userId}
	// consumerHandlernow.Init()

}
