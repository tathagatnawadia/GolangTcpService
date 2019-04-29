package Entities_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "../Entities"
	"../tests/mocks"
)

func TestClient_SendMessage(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock("Write")

	uid := 1
	text := "test message"
	expected := fmt.Sprintf("\nMessage from %d : %s\n", uid, text)

	client1 := &Client{uid, conn, make(chan RelayMessage), time.Now().String(), true, nil}
	client2 := &Client{uid, conn, make(chan RelayMessage), time.Now().String(), true, nil}
	message := CreateRelayMessage(text, "2", uid)

	go client1.ReceiveMessages()
	go client2.ReceiveMessages()
	client1.SendMessage(*message)
	time.Sleep(time.Second)

	assert.Equal(expected, conn.Result)
	conn.AssertExpectations(t)
}

func TestClient_AddToHistory(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock()
	uid := 1

	client1 := &Client{uid, conn, make(chan RelayMessage), time.Now().String(), true, nil}
	assert.Len(client1.History, 0)

	text := "test command 1"
	client1.AddToHistory(text)

	assert.Len(client1.History, 1)

	text = "test command 2"
	client1.AddToHistory(text)

	assert.Len(client1.History, 2)

	conn.AssertExpectations(t)
}

func TestClient_ReceiveMessages(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock("Write")
	uid := 1
	text := "test message"
	expected := fmt.Sprintf("\nMessage from %d : %s\n", uid, text)

	client1 := &Client{uid, conn, make(chan RelayMessage), time.Now().String(), true, nil}
	go client1.ReceiveMessages()
	message := CreateRelayMessage(text, "1", uid)

	client1.SendMessage(*message)
	time.Sleep(time.Second)

	assert.Equal(expected, conn.Result)

	conn.AssertExpectations(t)
}

func TestClient_GetUserId(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock()
	uid1 := 1

	client1 := &Client{
		uid1,
		conn,
		make(chan RelayMessage),
		time.Now().String(),
		true,
		nil,
	}
	assert.Equal(uid1, client1.GetUserId())
}

func TestClient_GetActive(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock()
	uid1 := 1
	uid2 := 2

	client1 := &Client{
		uid1,
		conn,
		make(chan RelayMessage),
		time.Now().String(),
		true,
		nil,
	}
	assert.True(client1.GetActive())

	client2 := &Client{
		uid2,
		conn,
		make(chan RelayMessage),
		time.Now().String(),
		false,
		nil,
	}
	assert.False(client2.GetActive())
}

func TestClient_SetActive(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock()
	uid1 := 1
	uid2 := 2

	client1 := &Client{
		uid1,
		conn,
		make(chan RelayMessage),
		time.Now().String(),
		true,
		nil,
	}
	assert.True(client1.GetActive())
	client1.SetActive(false)
	assert.False(client1.GetActive())


	client2 := &Client{
		uid2,
		conn,
		make(chan RelayMessage),
		time.Now().String(),
		false,
		nil,
	}
	assert.False(client2.GetActive())
	client2.SetActive(true)
	assert.True(client2.GetActive())
}
