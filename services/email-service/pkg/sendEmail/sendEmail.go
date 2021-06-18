package sendEmail

import (
	"crypto/tls"
	"fmt"
	"github.com/ankushChatterjee/jetpen/email-service/pkg/utils"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendEmail(emailContent string, subject string, from string, ownerMail string) error {
	mailUser := utils.GetEnvVar("MAIL_USERNAME")
	mailHost := utils.GetEnvVar("MAIL_HOST")
	mailPort := utils.GetEnvVar("MAIL_PORT")
	mailPassword := utils.GetEnvVar("MAIL_PASSWORD")

	e := email.NewEmail()
	e.From = from
	e.To = []string{ownerMail}
	e.Subject = subject
	e.HTML = []byte(emailContent)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         mailHost,
	}
	err := e.SendWithStartTLS(fmt.Sprintf("%s:%s", mailHost, mailPort), smtp.CRAMMD5Auth(mailUser, mailPassword), tlsconfig)
	if err != nil {
		return err
	}
	return nil
}
