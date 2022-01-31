package mails

import "github.com/stretchr/testify/mock"

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) Send(message *Message) error {
	args := m.Called(message)
	err := args.Error(0)
	return err
}
