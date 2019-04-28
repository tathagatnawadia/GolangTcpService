package mocks

import (
	"fmt"
	
	"github.com/stretchr/testify/mock"
)

type ConnMock struct {
	mock.Mock
	Result string
}

func (conn *ConnMock) Read(b []byte) (n int, err error) {
	conn.Called()

	return 0, nil
}

func (conn *ConnMock) Write(b []byte) (n int, err error) {
	conn.Called()

	fmt.Println("mocking net.Conn.Write()")
	conn.Result = string(b)

	return 0, nil
}

func NewConnMock(methods ...string) *ConnMock {
	connMock := new(ConnMock)
	for _, method := range methods {
		connMock.On(method).Return(nil)
	}

	return connMock
}
