package main

import (
	"fmt"
	producerHandlerv1 "yap-chat/producer/handler/v1"
	consumerHandlerv1 "yap-chat/consumer/handler/v1"
)

func main() {
	var userId int
	fmt.Print("Enter your user Id: ") 
	fmt.Scanln(&userId)
	producerHandlerv1.Init()
	consumerHandler := consumerHandlerv1.KafkaConsumer{
		UserID: userId,
	}
	consumerHandler.Init()
}
