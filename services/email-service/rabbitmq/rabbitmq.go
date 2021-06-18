package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/ankushChatterjee/jetpen/email-service/pkg/sendEmail"
	"github.com/ankushChatterjee/jetpen/email-service/pkg/utils"
	"github.com/streadway/amqp"
)

type EMailMessage struct {
	EmailContent string   `json:"emailContent"`
	Subject      string   `json:"subject"`
	From         string   `json:"from"`
	OwnerMail    string   `json:"ownerMail"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func StartConsumer() {
	conn, err := amqp.Dial(utils.GetEnvVar("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		utils.GetEnvVar("RABBITMQ_QUEUE_NAME"), // queue
		"",                                     // consumer
		true,                                   // auto-ack
		false,                                  // exclusive
		false,                                  // no-local
		false,                                  // no-wait
		nil,                                    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			msg := EMailMessage{}
			err = json.Unmarshal([]byte(d.Body), &msg)
			failOnError(err, "JSON Unmarshal failed")
			log.Printf("Received a message: %s", msg)
			err = sendEmail.SendEmail(msg.EmailContent, msg.Subject, msg.From, msg.OwnerMail)
			failOnError(err, "Email send failed")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
