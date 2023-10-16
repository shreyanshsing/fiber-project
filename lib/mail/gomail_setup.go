package gomail_lib

import (
	"fiber-project/config"
	"fiber-project/models"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

func prepareMail(sendersMail string, user *models.User) *gomail.Message {

	subject := "Confirmation Email"
	body := fmt.Sprintf("Thank you for registering to Fiber Project. Please confirm your email by clicking on following link - http://localhost:8080/api/auth/verify/%v", user.ID)

	// Create an email message
	mail := gomail.NewMessage()
	mail.SetHeader("From", sendersMail)
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/plain", body)

	return mail
}

func SendMail(user *models.User) {
	smtpServer := config.GetEnvoirmentVariable("SMTP_SERVER")
	smtpPort, _ := strconv.ParseUint(config.GetEnvoirmentVariable("SMTP_PORT"), 10, 64)
	senderEmail := config.GetEnvoirmentVariable("SENDER_EMAIL")
	senderPassword := config.GetEnvoirmentVariable("SENDER_PASSWORD")

	log.Println("sending verification email...", user.ID, senderEmail)

	mailBody := prepareMail(senderEmail, user)

	dailer := gomail.NewDialer(smtpServer, int(smtpPort), senderEmail, senderPassword)

	if err := dailer.DialAndSend(mailBody); err != nil {
		log.Println("failed to sent verification mail!")
		log.Panicln(err)
	}
	log.Println("Verification email sent successfully!")
}
