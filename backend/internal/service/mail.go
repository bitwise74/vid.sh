package service

import (
	"bitwise74/video-api/internal/model"
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject string, body string) error {
	if to == "" || subject == "" || body == "" {
		return errors.New("missing email parameters")
	}

	from := os.Getenv("MAIL_SENDER_ADDRESS")
	if to == from {
		return errors.New("invalid email address")
	}

	password := os.Getenv("MAIL_PASSWORD")

	m := gomail.NewMessage(func(m *gomail.Message) {
		m.SetHeader("From", from)
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)
	})

	smtpPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), smtpPort, from, password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendVerificationMail(t *model.Token, to string) error {
	from := os.Getenv("MAIL_SENDER_ADDRESS")
	if to == from {
		return errors.New("invalid email address")
	}

	password := os.Getenv("MAIL_PASSWORD")

	sslEnabled, err := strconv.ParseBool(os.Getenv("HOST_SSL_ENABLED"))
	if err != nil {
		sslEnabled = false
	}

	var s string
	if sslEnabled {
		s = "s"
	}

	verifLink := fmt.Sprintf("http%v://%v/verify?user_id=%v&token=%v",
		s, os.Getenv("HOST_DOMAIN"), t.UserID, t.Token)

	m := gomail.NewMessage(func(m *gomail.Message) {
		m.SetHeader("From", from)
		m.SetHeader("To", to)
		m.SetHeader("Subject", "Verify your email to start using vid.sh")
		m.SetBody("text/html", fmt.Sprintf("Click <a href='%v'>here</a> to verify your account.\n\nThis link will expire in 30 minutes", verifLink))
	})

	smtpPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), smtpPort, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendPasswdResetMail(t *model.Token, to string) error {
	from := os.Getenv("MAIL_SENDER_ADDRESS")
	if to == from {
		return errors.New("invalid email address")
	}

	password := os.Getenv("MAIL_PASSWORD")

	sslEnabled, err := strconv.ParseBool(os.Getenv("HOST_SSL_ENABLED"))
	if err != nil {
		sslEnabled = false
	}

	var s string
	if sslEnabled {
		s = "s"
	}

	verifLink := fmt.Sprintf("http%v://%v/reset-passwd?user_id=%v&token=%v",
		s, os.Getenv("HOST_DOMAIN"), t.UserID, t.Token)

	m := gomail.NewMessage(func(m *gomail.Message) {
		m.SetHeader("From", from)
		m.SetHeader("To", to)
		m.SetHeader("Subject", "Reset your password for vid.sh")
		m.SetBody("text/html", fmt.Sprintf("Click <a href='%v'>here</a> to reset your password.\n\nThis link will expire in 30 minutes. If you didn't request a password reset, please ignore this email.", verifLink))
	})

	smtpPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), smtpPort, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
