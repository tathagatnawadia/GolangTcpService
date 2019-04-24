package Entities
import (
	"net"
	"time"
	"sync"
	"strconv"
)
type Network struct {
	ClientList map[int]*Client
	ActiveConnections map[net.Conn]int
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
	r.ActiveConnections[conn] = client.User_id
	r.Mutex.Unlock()
	return client
}

func (r* Network) GetUserIdByConnection(conn net.Conn) (int, bool) {
	fetchedUserId, ok := r.ActiveConnections[conn]
	return fetchedUserId, ok
}

func (r* Network) GetClientById(userid int) (*Client, bool) {
	fetchedClient, ok := r.ClientList[userid]
	return fetchedClient, ok
}

func (r *Network) GetActiveClients(requestingClient *Client) string {
	var result = ""
	for index := range r.ClientList {
		if r.ClientList[index].Active == true && r.ClientList[index].User_id != requestingClient.User_id {
			result += strconv.Itoa(r.ClientList[index].User_id) + " "	
		}
	}

	if result == "" {
		result = "Looks lonely in here"
	}

	return result
}

func (r *Network) SendRelayMessage(message *RelayMessage, myClient *Client) (string, error) {
	for _, element := range message.ReceiptClients {
		if fetchedClient, ok := r.ClientList[element]; ok {
			if fetchedClient.User_id != myClient.User_id{
				go fetchedClient.SendMessage(RelayMessage{message.Message, message.From, nil})
			}
		}
	}
	return "ok", nil
}

func (r* Network) RemoveClientByConnection(conn net.Conn) bool {
	user_id := r.ActiveConnections[conn]
	client := r.ClientList[user_id]
	client.Active = false;
	delete(r.ActiveConnections, conn);
	conn.Close()
	return true
}


