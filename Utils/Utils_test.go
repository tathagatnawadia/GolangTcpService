package Utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"../tests/mocks"
	"../Utils"
)

func TestSendResponse(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock("Write")
	input := "Welcome to the hub !"
	expected := "Welcome to the hub !\n"

	Utils.SendResponse(input, conn)

	assert.NotEmpty(conn.Result)
	assert.Equal(expected, conn.Result)
	conn.AssertExpectations(t)
}

func TestSendPrompt(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock("Write")
	expected := ">> "

	Utils.SendPrompt(expected, conn)

	assert.NotEmpty(conn.Result)
	assert.Equal(expected, conn.Result)
	conn.AssertExpectations(t)
}

func TestSendBroadcast(t *testing.T) {
	assert := assert.New(t)
	conn := mocks.NewConnMock("Write")
	expected := "\nMessage from 5 : test message\n"

	Utils.SendBroadcast("test message", 5, conn)

	assert.NotEmpty(conn.Result)
	assert.Equal(expected, conn.Result)
	conn.AssertExpectations(t)
}

func TestPrintHelpText(t *testing.T) {
	assert := assert.New(t)
	expected := "----------HELP---------\n"
	expected += ">> LIST\n"
	expected += ">> IDENTIFY\n"
	expected += ">> RELAY #How you doing user 1,3 and 4 ?? #1,3,4\n"
	expected += ">> EXIT\n"
	expected += "-----------------------\n"
	conn := mocks.NewConnMock("Write")

	Utils.PrintHelpText(conn)
	assert.NotEmpty(conn.Result)
	assert.Equal(expected, conn.Result)
	conn.AssertExpectations(t)
}
