package mail

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/stepnovij/authBoilerplate/utils"
)

var username = os.Getenv("SMTP_USERNAME")
var password = os.Getenv("SMTP_PASSWORD")
var host = os.Getenv("SMTP_HOST")
var port = os.Getenv("SMTP_PORT")
var from = os.Getenv("EMAIL_FROM")

func sendMail(to []string, from string, msg []byte) {
	// Set up authentication information
	addr := utils.Concatanete(host, port)
	auth := smtp.PlainAuth("", username, password, host)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	//msg := []byte("To: stepnovij@gmail.com\r\n" +
	//	"Subject: discount Gophers!\r\n" +
	//	"\r\n" +
	//	"This is the email body.\r\n")

	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		panic(err)
	}
}

func SendConfirmationEmail(to []string, confirmationLink string) {
	msgText := fmt.Sprintf("Subject: Email confirmation\r\n"+
		"Thank you for registration. "+
		"Please confirm your email: %v", confirmationLink)
	msg := []byte(msgText)
	sendMail(to, from, msg)
}

func SendReferralLinkEmail(to []string, referralLink string) {
	msgText := fmt.Sprintf("Subject: Referral link\r\n"+
		"Thank you for confirmation your mail."+
		"This is your referral link: %v", referralLink)
	msg := []byte(msgText)
	sendMail(to, from, msg)
}
