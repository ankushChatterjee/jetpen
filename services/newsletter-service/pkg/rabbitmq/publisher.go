package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/utils"
	"github.com/streadway/amqp"
)

var ch *amqp.Channel

const EmailFrom = "JetPen <no-reply@jetpen.com>"
const JunkMail = "target@jetpen.com"

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("ERROR :%s: %s", msg, err)
	}
}

func fatalOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("ERROR : %s: %s", msg, err)
	}
}

func Init() {
	conn, err := amqp.Dial(utils.GetEnvVar("RABBITMQ_URL"))
	fatalOnError(err, "Failed to connect to RabbitMQ")
	ch, err = conn.Channel()
	fatalOnError(err, "Failed to open a channel")
}

func PublishLetter(to []string, subject string, message string) error {
	emailData := map[string]interface{}{
		"emailContent": message,
		"subject":      subject,
		"from":         EmailFrom,
		"ownerMail":    JunkMail,
		"sendTo":       to,
	}
	data, _ := json.Marshal(emailData)
	err := ch.Publish(
		"",                                     // exchange
		utils.GetEnvVar("RABBITMQ_QUEUE_NAME"), // routing key
		false,                                  // mandatory
		false,                                  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        data,
		})
	failOnError(err, "Failed to publish a message")
	return err
}
