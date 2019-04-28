package Entities

import (
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

const ActiveClientsEmptyListResult = "Looks lonely in here"

type Network struct {
	ClientList        map[int]IClient
	ActiveConnections map[net.Conn]int
	Total             int
	Mutex             sync.Mutex
}

func (r *Network) assignAddressToNode() int {
	r.Total = r.Total + 1
	return r.Total
}

func (r *Network) Register(conn net.Conn) *Client {
	r.Mutex.Lock()
	client := &Client{r.assignAddressToNode(), conn, make(chan RelayMessage), time.Now().String(), true, nil}
	r.ClientList[client.GetUserId()] = client
	r.ActiveConnections[conn] = client.GetUserId()
	r.Mutex.Unlock()
	return client
}

func (r *Network) GetUserIdByConnection(conn net.Conn) (int, bool) {
	fetchedUserId, ok := r.ActiveConnections[conn]
	return fetchedUserId, ok
}

func (r *Network) GetClientById(userid int) (IClient, bool) {
	fetchedClient, ok := r.ClientList[userid]
	return fetchedClient, ok
}

func (r *Network) GetActiveClients(requestingClient IClient) string {
	var result = ""
	keys := make([]int, 0)
	for index := range r.ClientList {
		keys = append(keys, index)
	}
	sort.Ints(keys)

	for _, index := range keys {
		if r.ClientList[index].GetActive() == true && r.ClientList[index].GetUserId() != requestingClient.GetUserId() {
			result += strconv.Itoa(r.ClientList[index].GetUserId()) + " "
		}
	}

	if result == "" {
		result = ActiveClientsEmptyListResult
	}

	return result
}

func (r *Network) SendRelayMessage(message *RelayMessage, myClient IClient) (string, error) {
	for _, element := range message.ReceiptClients {
		if fetchedClient, ok := r.ClientList[element]; ok {
			if fetchedClient.GetUserId() != myClient.GetUserId() {
				//@todo: what if the user_id is valid but not active, we are sending messages to dead clients
				go fetchedClient.SendMessage(RelayMessage{message.Message, message.From, nil})
			}
		}
	}
	//@todo: quite pointless right now
	return "ok", nil
}

func (r *Network) RemoveClientByConnection(conn net.Conn) bool {
	user_id, ok := r.ActiveConnections[conn]
	//@todo: handle remove client false 
	if !ok {
		return false
	}
	client := r.ClientList[user_id]
	client.SetActive(false)
	delete(r.ActiveConnections, conn)
	conn.Close()
	//@todo: quite pointless to make permanent ok response
	return true
}

func NewNetwork() (*Network, error) {
	return &Network{make(map[int]IClient), make(map[net.Conn]int), 0, sync.Mutex{}}, nil
}
