package Entities
import (
	"net"
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