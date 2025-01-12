package main

import (
	"fmt"
	producerHandlerv1 "yap-chat/producer/handler/v1"
	consuerHandlerv1 "yap-chat/consumer/handler/v1"
)

func main() {
	var userId int
	fmt.Print("Enter your user Id: ") 
	fmt.Scanln(&userId)
	producerHandlerv1.Init()
	consuerHandlerv1.Init()
}