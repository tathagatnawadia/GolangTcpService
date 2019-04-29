package mocks

import (
	"fmt"
	"net"
	"time"
	
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

func (conn *ConnMock) Close() error {
	conn.Called()

	return nil
}

func (conn *ConnMock) LocalAddr() net.Addr {
	conn.Called()

	return nil
}

func (conn *ConnMock) RemoteAddr() net.Addr {
	conn.Called()

	return nil
}


func (conn *ConnMock) SetDeadline(t time.Time) error {
	conn.Called()

	return nil
}

func (conn *ConnMock) SetReadDeadline(t time.Time) error {
	conn.Called()

	return nil
}

func (conn *ConnMock) SetWriteDeadline(t time.Time) error {
	conn.Called()

	return nil
}

func NewConnMock(methods ...string) *ConnMock {
	connMock := new(ConnMock)
	for _, method := range methods {
		connMock.On(method).Return(nil)
	}

	return connMock
}
