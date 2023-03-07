package email

import (
	"fmt"
	"log"
	"net/smtp"
)

var config Config

func SetConfig(cfg Config) {
	config = cfg
}

func SendMail(to, subject, body string) {
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n\r\n"+
		"%s\r\n", config.From, to, subject, body))

	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err := smtp.SendMail(addr, auth, config.From, []string{to}, msg)

	if err != nil {
		log.Fatal(err)
	}
}
