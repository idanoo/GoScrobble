package goscrobble

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendEmail(destName string, destEmail string, subject string, content string) error {
	from := mail.NewEmail(os.Getenv("MAIL_FROM_NAME"), os.Getenv("MAIL_FROM_ADDRESS"))
	to := mail.NewEmail(destName, destEmail)
	message := mail.NewSingleEmail(from, subject, to, content, "")
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	_, err := client.Send(message)
	return err
}
