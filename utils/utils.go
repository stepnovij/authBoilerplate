package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/badoux/checkmail"
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

func EncryptPassword(password string) []byte {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hash to store:", string(hash))
	return hash
}

func comparePasswords(password string, hash []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		panic(err)
	}
	return true
}

func ValidatEmail(email string) error {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return err
	}
	return nil
}

func Concatanete(host string, port string) string {
	var str bytes.Buffer
	delimiter := ":"

	list := []string{host, delimiter, port}

	for _, l := range list {
		str.WriteString(l)
	}
	return str.String()
}
