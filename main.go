package main

import (
	"fmt"
	"sync"
	consumerHandlerv1 "yap-chat/consumer/handler/v1"
	producerHandlerv1 "yap-chat/producer/handler/v1"
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

	var wg sync.WaitGroup
	wg.Add(1) // Add a WaitGroup to block the main function
	wg.Wait()
}
