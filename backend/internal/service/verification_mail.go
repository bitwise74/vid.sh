package service

import (
	"bitwise74/video-api/internal/model"
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendVerificationMail(t *model.VerificationToken, sendTo string) error {
	from := os.Getenv("MAIL_SENDER_ADDRESS")
	if sendTo == from {
		return errors.New("invalid email address")
	}

	password := os.Getenv("MAIL_PASSWORD")

	m := gomail.NewMessage()

	sslEnabled, err := strconv.ParseBool(os.Getenv("HOST_SSL_ENABLE"))
	if err != nil {
		sslEnabled = false
	}

	var s string
	if sslEnabled {
		s = "s"
	}

	verifLink := fmt.Sprintf("http%v://%v/verify?user_id=%v&token=%v",
		s, os.Getenv("HOST_DOMAIN"), t.UserID, t.Token)

	m.SetHeader("From", from)
	m.SetHeader("To", sendTo)
	m.SetHeader("Subject", "Verify your email to start using vid.sh")
	m.SetBody("text/html", fmt.Sprintf("Click <a href='%v'>here</a> to verify your account.\n\nThis link will expire in 30 minutes", verifLink))

	smtpPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), smtpPort, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
