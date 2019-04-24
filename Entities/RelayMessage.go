package Entities

import (
	"strings"
	"strconv"
	"unsafe"
)

type RelayMessage struct {
	Message string
	From int
	ReceiptClients []int
}

func (m *RelayMessage) ValidateMessageLength(maxSize int) bool {
	return (len(m.Message) + int(unsafe.Sizeof(m.Message)))/1000 <= maxSize
}

func (m *RelayMessage) ValidateRecieverCount(maxSize int) bool {
	return len(m.ReceiptClients) <= maxSize
}

func CreateRelayMessage(message string, recievers string, from int) *RelayMessage {
	relayMessage := &RelayMessage{message, from, nil}
	receipts := strings.Split(recievers, ",")

	for index := range receipts {
		user_id, err := strconv.Atoi(receipts[index])
		if err == nil {
	 		relayMessage.ReceiptClients = append(relayMessage.ReceiptClients, user_id)
		}
	}

	return relayMessage
}