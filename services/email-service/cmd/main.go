package main

import (
	"log"

	"github.com/ankushChatterjee/jetpen/email-service/rabbitmq"
)

func main() {
	log.Println("[INFO] Staring email-service")
	rabbitmq.StartConsumer()
}
