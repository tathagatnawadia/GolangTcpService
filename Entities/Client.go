package Entities
import (
	"net"
	"../Utils"
)

type Client struct {
	User_id     int
	Handler net.Conn
	Incoming chan RelayMessage
	Timestamp string
	Active bool
	History []string
}

func (r *Client) SendMessage(myMessage RelayMessage){
	r.Incoming <- myMessage
}

func (r *Client) AddToHistory(command string){
	r.History = append(r.History, command)
}

func (r *Client) ReceiveMessages(){
	for {
		messagePacket := <-r.Incoming
		Utils.SendBroadcast(messagePacket.Message, messagePacket.From, r.Handler)
	}
}
