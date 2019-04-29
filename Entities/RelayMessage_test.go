package Entities_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "relay_solution/Entities"
)

func TestRelayMessage_ValidateMessageLength(t *testing.T) {
	assert := assert.New(t)

	uid := 1
	toUid := "2,3"
	text := "test message"

	message := CreateRelayMessage(text, toUid, uid)

	assert.True(message.ValidateMessageLength(100))
	assert.False(message.ValidateMessageLength(-1))
}

func TestRelayMessage_ValidateRecieverCount(t *testing.T) {
	assert := assert.New(t)

	uid := 1
	toUid := "2,3,5,11"
	text := "test message"

	message := CreateRelayMessage(text, toUid, uid)

	assert.True(message.ValidateRecieverCount(100))
	assert.False(message.ValidateRecieverCount(1))
}

func TestCreateRelayMessage(t *testing.T) {
	assert := assert.New(t)

	uid := 1
	toUid := "2,3"
	text := "test message"

	message := CreateRelayMessage(text, toUid, uid)
	assert.IsType(&RelayMessage{}, message)
	assert.Equal(text, message.Message)
	assert.Equal(uid, message.From)
	assert.Equal([]int{2,3}, message.ReceiptClients)
}
