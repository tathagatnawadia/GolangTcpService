package Utils

import (
	"net"
)

func SendResponse(res string, conn net.Conn) {
    conn.Write([]byte(res+"\n"))
}

func SendPrompt(res string, conn net.Conn) {
    conn.Write([]byte(res))
}