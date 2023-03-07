package util

import (
	"math/rand"
	"netradio/internal/model"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const VerificationCodeTTL = time.Minute * 10

func IsPasswordValid(password string) bool {
	return regexp.MustCompile(`.{1,30}`).MatchString(password)
}

func IsEmailValid(email string) bool {
	return regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenVerificationCode(email string) model.VerificationCode {
	b := make([]rune, 6)
	for i := range b {
		b[i] = rune(rand.Intn(10) + int('0'))
	}

	var result model.VerificationCode
	result.Email = email
	result.Value = string(b)
	result.Expires = time.Now().Add(VerificationCodeTTL)

	return result
}
