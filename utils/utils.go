package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io"
	"regexp"
)

func HashPassword(password string) string {
	bts, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(bts)
}

func GinContextFromContext(ctx context.Context) *gin.Context {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		panic(err)
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		panic(err)
	}

	return gc
}

func Encrypt(key []byte, message string) (encoded string, err error) {
	plainText := []byte(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.RawStdEncoding.EncodeToString(cipherText), err
}

func Decrypt(key []byte, secure string) (decoded string, err error) {
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = fmt.Errorf("ciphertext block size is too short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}

func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func IsPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}

	specialCharactersRegex := regexp.MustCompile(`\W`)
	lettersRegex := regexp.MustCompile(`[a-z]`)
	numbersRegex := regexp.MustCompile(`[0-9]`)
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	lowercaseRegex := regexp.MustCompile(`[a-z]`)

	if !specialCharactersRegex.MatchString(password) {
		return false
	} else if !lettersRegex.MatchString(password) {
		return false
	} else if !numbersRegex.MatchString(password) {
		return false
	} else if !uppercaseRegex.MatchString(password) {
		return false
	} else if !lowercaseRegex.MatchString(password) {
		return false
	}

	return true
}
