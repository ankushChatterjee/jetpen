package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/utils"
	"github.com/streadway/amqp"
)

var ch *amqp.Channel

const EmailFrom = "JetPen <no-reply@jetpen.com>"
const JunkMail = "target@jetpen.com"
const EmailStyle = "<link rel=\"stylesheet\" href=\"https://unpkg.com/purecss@2.0.6/build/pure-min.css\" integrity=\"sha384-Uu6IeWbM+gzNVXJcM9XV3SohHtmWE+3VGi496jvgX1jyvDTXfdK+rfZc8C1Aehk5\" crossorigin=\"anonymous\">"
const EmaillFooter = "To Unsubscribe please <a href=\"%s\">Click Here</a>"

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

func PublishLetter(to map[string]string, subject string, message string, nid string) error {
	frontendHost := utils.GetEnvVar("FRONTEND_HOST")

	for email, token := range to {
		completeMessage := EmailStyle + message + fmt.Sprintf(EmaillFooter, frontendHost + "/unsub/?email=" + email + "&token=" + token + "&nid=" + nid)
		emailData := map[string]interface{}{
			"emailContent": completeMessage,
			"subject":      subject,
			"from":         EmailFrom,
			"ownerMail":    email,
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
		if err != nil {
			return err
		}
	}
	return nil
}
