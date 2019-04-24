package Utils

import (
	"net"
	"strconv"
)

func SendResponse(res string, conn net.Conn) {
    conn.Write([]byte(res+"\n"))
}

func SendPrompt(res string, conn net.Conn) {
    conn.Write([]byte(res))
}

func SendBroadcast(res string,from int ,conn net.Conn) {
	res = "\n Message from " + strconv.Itoa(from) + " : " + res + "\n"
	conn.Write([]byte(res))
}