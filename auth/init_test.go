package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitService(t *testing.T) {
	InitService()
	assert.NotNil(t, Service)
}
