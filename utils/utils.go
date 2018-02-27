package utils

import (
	"math/rand"
	"time"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"github.com/badoux/checkmail"
	"net/smtp"
	"log"
	"bytes"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


func EncryptPassword(password string) ([]byte) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hash to store:", string(hash))
	return hash
}


func comparePasswords(password string, hash []byte) (bool){
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		panic(err)
	}
	return true
}


func ValidatEmail(email string) (error){
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return err
	}
	return nil
}


func sendEmail(sender_email string, recipient_email string, text string) (){
	c, err := smtp.Dial("mail.example.com:25")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Set the sender and recipient.
	c.Mail(sender_email)
	c.Rcpt(recipient_email)

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(text)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
	return
}


func SendConfirmationEmail(recipient_email string, hash string) {
	sender_email := "email"
	text := "Thank you for registration. Please confirm your email"
	// text + hash?
	sendEmail(sender_email, recipient_email, text)
}


func SendReferralLinkEmail(recipient_email string, referralLink string) {
	sender_email := "email"
	text := "Here is your referral link."
	// text + hash?
	sendEmail(sender_email, recipient_email, text)
}