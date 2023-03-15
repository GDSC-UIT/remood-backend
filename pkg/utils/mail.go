package utils

import (
	"crypto/tls"	

	gomail "gopkg.in/mail.v2"
)

func SendMail(message string, toList []string) error {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "remoodgdsc@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", toList[0])

	// Set E-Mail subject
	m.SetHeader("Subject", "REMOOD: Reseting password email")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", message)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "remoodgdsc@gmail.com", "remood@123")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
