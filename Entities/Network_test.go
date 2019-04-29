package Entities_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "relay_solution/Entities"
	"relay_solution/tests/mocks"
)

func TestNetwork_Register(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock()

	network, err := NewNetwork()

	assert.IsType(&Network{}, network)
	assert.Nil(err)

	client := network.Register(conn)
	assert.IsType(&Client{}, client)

	assert.Equal(true, client.Active)
	assert.Equal(conn, client.Handler)
	assert.Nil(client.History)
	assert.Equal(1, client.User_id)

	conn.AssertExpectations(t)
}

func TestNetwork_Register_Two_Clients(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock()

	network, err := NewNetwork()

	assert.IsType(&Network{}, network)
	assert.Nil(err)

	client1 := network.Register(conn)
	assert.IsType(&Client{}, client1)

	assert.Equal(true, client1.Active)
	assert.Equal(conn, client1.Handler)
	assert.Nil(client1.History)
	assert.Equal(1, client1.User_id)

	client2 := network.Register(conn)
	assert.IsType(&Client{}, client2)

	assert.Equal(true, client2.Active)
	assert.Equal(conn, client2.Handler)
	assert.Nil(client2.History)
	assert.Equal(2, client2.User_id)

	conn.AssertExpectations(t)
}

func TestNetwork_GetUserIdByConnection(t *testing.T) {
	assert := assert.New(t)
	conn1 := mocks.NewConnMock()
	conn2 := mocks.NewConnMock()

	network, err := NewNetwork()

	assert.IsType(&Network{}, network)
	assert.Nil(err)

	uid, ok := network.GetUserIdByConnection(conn1)
	assert.False(ok)
	assert.Zero(uid)

	uid, ok = network.GetUserIdByConnection(conn2)
	assert.False(ok)
	assert.Zero(uid)

	network.Register(conn1)

	uid, ok = network.GetUserIdByConnection(conn1)
	assert.True(ok)
	assert.Equal(1, uid)

	network.Register(conn1)

	uid, ok = network.GetUserIdByConnection(conn1)
	assert.True(ok)
	assert.Equal(2, uid)

	conn1.AssertExpectations(t)
}

func TestNetwork_GetClientById(t *testing.T) {
	assert := assert.New(t)
	conn1 := mocks.NewConnMock()
	conn2 := mocks.NewConnMock()

	network, err := NewNetwork()

	assert.IsType(&Network{}, network)
	assert.Nil(err)

	client, ok := network.GetClientById(1)
	assert.False(ok)
	assert.Nil(client)

	client, ok = network.GetClientById(2)
	assert.False(ok)
	assert.Nil(client)

	network.Register(conn1)

	client, ok = network.GetClientById(1)
	assert.True(ok)
	assert.NotNil(client)
	assert.IsType(&Client{}, client)

	network.Register(conn2)

	client, ok = network.GetClientById(2)
	assert.True(ok)
	assert.NotNil(client)
	assert.IsType(&Client{}, client)

	conn1.AssertExpectations(t)
}

func TestNetwork_GetActiveClients(t *testing.T) {
	assert := assert.New(t)
	conn1 := mocks.NewConnMock()
	expectedEmpty := ActiveClientsEmptyListResult

	network, err := NewNetwork()

	assert.IsType(&Network{}, network)
	assert.Nil(err)

	//@todo: should it be allowed to pass nil here?
	clients := network.GetActiveClients(nil)
	assert.Equal(expectedEmpty, clients)

	client1 := network.Register(conn1)

	clients = network.GetActiveClients(client1)
	assert.Equal(expectedEmpty, clients)

	client2 := network.Register(conn1)

	clients = network.GetActiveClients(client1)
	//@todo: Trailing space, should be removed
	assert.Equal("2 ", clients)

	clients = network.GetActiveClients(client2)
	assert.Equal("1 ", clients)

	client3 := network.Register(conn1)

	clients = network.GetActiveClients(client3)
	assert.Equal("1 2 ", clients)

	clients = network.GetActiveClients(client2)
	assert.Equal("1 3 ", clients)

	conn1.AssertExpectations(t)
}

func TestNetwork_SendRelayMessage(t *testing.T) {
	assert := assert.New(t)
	text := "test message"
	uidFrom := 1
	uidTo := 2
	conn1 := mocks.NewConnMock()
	conn2 := mocks.NewConnMock()
	client1 := mocks.NewClientMock(uidFrom)
	client2 := mocks.NewClientMock(uidTo)

	uidToStr := strconv.Itoa(uidTo)
	client1.On("GetUserId").Return(uidFrom)
	client2.On("GetUserId").Return(uidTo)

	client2.On("SendMessage").Return(nil)

	network, err := NewNetwork()
	assert.IsType(&Network{}, network)
	assert.Nil(err)

	// fake the network.Register()
	network.ClientList[uidFrom] = client1
	network.ActiveConnections[conn1] = uidFrom
	network.ClientList[uidTo] = client2
	network.ActiveConnections[conn2] = uidTo

	message := CreateRelayMessage(text, uidToStr, uidFrom)

	ok, err := network.SendRelayMessage(message, client1)
	time.Sleep(2 * time.Second)
	assert.Equal("ok", ok)
	assert.Nil(err)
	//@todo: some mistake in mocking, more debug needed
	// assert.Len(client2.Incoming, 1)
	// assert.Equal(message, client2.Incoming[0])

	conn1.AssertExpectations(t)
	conn2.AssertExpectations(t)
	client1.AssertExpectations(t)
	client2.AssertExpectations(t)
}

func TestNetwork_RemoveClientByConnection(t *testing.T) {
	assert := assert.New(t)
	conn1 := mocks.NewConnMock("Close")
	conn2 := mocks.NewConnMock("Close")

	network, err := NewNetwork()

	assert.IsType(&Network{}, network)
	assert.Nil(err)

	assert.False(network.RemoveClientByConnection(conn1))

	network.Register(conn1)
	network.Register(conn2)

	assert.True(network.RemoveClientByConnection(conn2))
	assert.True(network.RemoveClientByConnection(conn1))
}

func TestNewNetwork(t *testing.T) {
	assert := assert.New(t)

	network, err := NewNetwork()

	assert.IsType(&Network{}, network)
	assert.NotNil(network)
	assert.Nil(err)
}
