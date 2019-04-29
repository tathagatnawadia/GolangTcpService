package mocks

import (
	"net"

	"github.com/stretchr/testify/mock"

	. "relay_solution/Entities"
	"relay_solution/Utils"
)

type ClientMock struct {
	mock.Mock
	User_id   int
	Handler   net.Conn
	// Some debug required on how to test this - differs from channel
	Incoming  []*RelayMessage
	Timestamp string
	Active    bool
	History   []string
}

func (r *ClientMock) GetUserId() int {
	r.Called()

	return r.User_id
}

func (r *ClientMock) GetActive() bool {
	r.Called()

	return r.Active
}

func (r *ClientMock) SetActive(v bool) {
	r.Called()
	r.Active = v
}

func (r *ClientMock) SendMessage(myMessage RelayMessage) {
	r.Called()
	r.Incoming = append(r.Incoming, &myMessage)
}

func (r *ClientMock) AddToHistory(command string) {
	r.Called()
	r.History = append(r.History, command)
}

func (r *ClientMock) ReceiveMessages() {
	r.Called()
	var messagePacket RelayMessage
	messagePacket, r.Incoming = *r.Incoming[len(r.Incoming)-1], r.Incoming[:len(r.Incoming)-1]
	Utils.SendBroadcast(messagePacket.Message, messagePacket.From, r.Handler)
}

// func NewClientMock(methods ...string) *ClientMock {
// 	clientMock := new(ClientMock)
// 	for _, method := range methods {
// 		clientMock.On(method).Return(nil)
// 	}

// 	return clientMock
// }

func NewClientMock(fakeUserId int) *ClientMock {
	clientMock := new(ClientMock)
	clientMock.User_id = fakeUserId

	return clientMock
}
