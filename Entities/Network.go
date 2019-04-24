package Entities
import (
	"net"
	"time"
)
type Network struct {
	ClientList map[int]*Client
	Total int
}

func (r *Network) AssignAddressToNode() int {
	r.Total = r.Total + 1
	return r.Total
}
 
func (r* Network) AddNewClient(conn net.Conn) *Client {
	client := &Client{r.AssignAddressToNode(), conn, make(chan RelayMessage), time.Now().String(), true, nil}
	r.ClientList[client.User_id] = client
	return client
}

