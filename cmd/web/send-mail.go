package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/prranavv/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorlog.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email-template/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		msgTosend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgTosend)
	}
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	}
	log.Println("email sent")
}
