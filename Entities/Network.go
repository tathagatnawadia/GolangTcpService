package Entities
import (
	"net"
	"time"
	"sync"
)
type Network struct {
	ClientList map[int]*Client
	Total int
	Mutex sync.Mutex
}

func (r *Network) AssignAddressToNode() int {
	r.Total = r.Total + 1
	return r.Total
}
 
func (r* Network) AddNewClient(conn net.Conn) *Client {
	r.Mutex.Lock()
	client := &Client{r.AssignAddressToNode(), conn, make(chan RelayMessage), time.Now().String(), true, nil}
	r.ClientList[client.User_id] = client
	r.Mutex.Unlock()
	return client
}

