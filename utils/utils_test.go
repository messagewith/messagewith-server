package utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http/httptest"
	"testing"
)

func TestHashPassword(t *testing.T) {
	firstHash := HashPassword("hello world")
	err := bcrypt.CompareHashAndPassword([]byte(firstHash), []byte("hello world"))

	assert.Nil(t, err)
}

func TestEncryptAndDecrypt(t *testing.T) {
	key := []byte("")
	_, err := Encrypt(key, "hello world")
	assert.NotNil(t, err)

	key = []byte("%123")
	_, err = Encrypt(key, "hello world")
	assert.NotNil(t, err)

	key = []byte("%123psu689")
	_, err = Encrypt(key, "hello world")
	assert.NotNil(t, err)

	key = []byte("%123psu689123456790123p756asdasd][][[][sa")
	_, err = Encrypt(key, "hello world")
	assert.NotNil(t, err)

	key = []byte("")
	_, err = Decrypt(key, "hello world")
	assert.NotNil(t, err)

	key = []byte("%123")
	_, err = Decrypt(key, "hello world")
	assert.NotNil(t, err)

	key = []byte("%123psu689")
	_, err = Decrypt(key, "hello world")
	assert.NotNil(t, err)

	key = []byte("%123psu689123456790123p756asdasd][][[][sa")
	_, err = Decrypt(key, "hello world")
	assert.NotNil(t, err)

	key = []byte("%123psu689123456790123p756sy%221")
	encryptedValue, err := Encrypt(key, "hello world")
	assert.Nil(t, err)

	decryptedValue, err := Decrypt(key, encryptedValue)
	assert.Nil(t, err)

	assert.Equal(t, decryptedValue, "hello world")

	encryptedValue, err = Encrypt(key, "world")
	assert.Nil(t, err)
	decryptedValue, err = Decrypt(key, encryptedValue)
	assert.Nil(t, err)

	assert.Equal(t, decryptedValue, "world")
}

func TestGinContextFromContext(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/", func(ginCtx *gin.Context) {
		ctx := context.WithValue(ginCtx.Request.Context(), "GinContextKey", ginCtx)
		testCtx := GinContextFromContext(ctx)
		assert.Equal(t, ginCtx, testCtx)
	})

	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
}
