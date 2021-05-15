package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
	"github.com/streadway/amqp"
)

var ch *amqp.Channel

const VERIFICATION_EMAIL_SUBJECT = "Welcome to JetPen, Verify your email address."
const VERIFICATION_EMAIL_CONTENT = "<h1>Welcome to Jetpen</h1><p>Go to <a href=\"%s\">%s</a></p>"
const VERIFICATION_EMAIL_FROM = "JetPen <contact@jetpen.com>"

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
	url := utils.GetEnvVar("BASE_URL") + "/verify-email/?token=" + token + "&username=" + username
	emaiData := map[string]string{
		"emailContent": fmt.Sprintf(VERIFICATION_EMAIL_CONTENT, url, url),
		"subject":      VERIFICATION_EMAIL_SUBJECT,
		"from":         VERIFICATION_EMAIL_FROM,
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
