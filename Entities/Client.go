package Entities

import (
	"net"
	"../Utils"
)

type IClient interface {
	GetUserId() int
	GetActive() bool
	SetActive(bool)

	SendMessage(myMessage RelayMessage)
	AddToHistory(command string)
	ReceiveMessages()
}

//@todo: not a good to expose properties as public, should be private with getters and setters if needed - following SOLID
type Client struct {
	User_id   int
	Handler   net.Conn
	Incoming  chan RelayMessage
	Timestamp string
	Active    bool
	History   []string
}

func (r *Client) GetUserId() int {
	return r.User_id
}

func (r *Client) GetActive() bool {
	return r.Active
}

func (r *Client) SetActive(v bool) {
	r.Active = v
}

func (r *Client) SendMessage(myMessage RelayMessage) {
	r.Incoming <- myMessage
}

func (r *Client) AddToHistory(command string) {
	r.History = append(r.History, command)
}

func (r *Client) ReceiveMessages() {
	for {
		messagePacket := <-r.Incoming
		Utils.SendBroadcast(messagePacket.Message, messagePacket.From, r.Handler)
	}
}
