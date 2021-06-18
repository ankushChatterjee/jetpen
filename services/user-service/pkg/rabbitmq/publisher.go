package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
	"github.com/streadway/amqp"
)

var ch *amqp.Channel

const VerificationEmailSubject = "Welcome to JetPen, Verify your email address."
const VerificationEmailContent = "<h1>Welcome to Jetpen</h1><p>Go to <a href=\"%s\">%s</a></p>"
const VerificationEmailFrom = "JetPen <no-reply@jetpen.com>"

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

func PublishVerificationEmail(to string, token string, username string) error {
	frontendHost := utils.GetEnvVar("FRONTEND_HOST")
	url := utils.GetEnvVar("BASE_URL") + frontendHost + "/verify-email/?token=" + token + "&username=" + username
	emaiData := map[string]string{
		"emailContent": fmt.Sprintf(VerificationEmailContent, url, url),
		"subject":      VerificationEmailSubject,
		"from":         VerificationEmailFrom,
		"ownerMail":    to,
	}
	data, _ := json.Marshal(emaiData)
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
