package internal

import (
	"os"

	"gopkg.in/gomail.v2"
)

// Mailer Mailer
type Mailer struct {
	SMTPHost string
	SMTPPort int
	Account  string
	Password string
}

// Init Init
func (mailer *Mailer) Init(smtpHost string, smtpPort int, account, password string) {
	mailer.SMTPHost = smtpHost
	mailer.SMTPPort = smtpPort
	mailer.Account = account
	mailer.Password = password
}

// Send Send
func (mailer *Mailer) Send(receivers []string, subject, body string, attachments *[]string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailer.Account)
	m.SetHeader("To", receivers...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if attachments != nil {
		for _, attachment := range *attachments {
			if _, err := os.Stat(attachment); err == nil {
				m.Attach(attachment)
			}
		}
	}
	d := gomail.NewDialer(mailer.SMTPHost, mailer.SMTPPort, mailer.Account, mailer.Password)
	return d.DialAndSend(m)
}
