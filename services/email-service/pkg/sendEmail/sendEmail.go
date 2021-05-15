package sendEmail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/ankushChatterjee/jetpen/email-service/pkg/utils"
	"github.com/jordan-wright/email"
)

func SendEmail(emailContent string, sendTo []string, subject string, from string, ownerMail string) error {
	mailUser := utils.GetEnvVar("MAIL_USERNAME")
	mailHost := utils.GetEnvVar("MAIL_HOST")
	mailPort := utils.GetEnvVar("MAIL_PORT")
	mailPassword := utils.GetEnvVar("MAIL_PASSWORD")
	mailMaxSend, err := strconv.Atoi(utils.GetEnvVar("MAIL_MAX_SEND"))
	if err != nil {
		return errors.New("Invalid MAIL_MAX_SEND")
	}

	min := 0
	max := mailMaxSend - 1

	if len(sendTo) == 0 {
		e := email.NewEmail()
		e.From = from
		e.To = []string{ownerMail}
		e.Subject = subject
		e.HTML = []byte(emailContent)
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         mailHost,
		}
		err = e.SendWithStartTLS(fmt.Sprintf("%s:%s", mailHost, mailPort), smtp.CRAMMD5Auth(mailUser, mailPassword), tlsconfig)
		if err != nil {
			return err
		}
		return nil
	}

	for min < len(sendTo) {
		if max >= len(sendTo)-1 {
			max = len(sendTo) - 1
		}
		e := email.NewEmail()
		e.From = from
		e.To = []string{ownerMail}
		e.Bcc = sendTo[min : max+1]
		e.Subject = subject
		e.HTML = []byte(emailContent)
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         mailHost,
		}

		err = e.SendWithStartTLS(fmt.Sprintf("%s:%s", mailHost, mailPort), smtp.CRAMMD5Auth(mailUser, mailPassword), tlsconfig)
		if err != nil {
			return err
		}
		min = max + 1
		max = min + mailMaxSend - 1
	}
	return nil
}
