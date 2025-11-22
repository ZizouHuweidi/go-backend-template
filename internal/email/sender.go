package email

import (
	"fmt"
	"net/smtp"
)

type Sender struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

func NewSender(host string, port int, username, password, sender string) *Sender {
	return &Sender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Sender:   sender,
	}
}

func (s *Sender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.Sender, to, subject, body))

	if s.Username != "" && s.Password != "" {
		auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
		return smtp.SendMail(addr, auth, s.Sender, []string{to}, msg)
	}

	// No auth (e.g. MailHog)
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Mail(s.Sender); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}
