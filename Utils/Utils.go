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
	res = "\nMessage from " + strconv.Itoa(from) + " : " + res + "\n"
	conn.Write([]byte(res))
}

func PrintHelpText(conn net.Conn) {
	res := "----------HELP---------\n"
	res += "> LIST\n"
	res += "> IDENTIFY\n"
	res += "> RELAY #How you doing user 1,3 and 4 ?? #1,3,4\n"
	res += "-----------------------\n"
    conn.Write([]byte(res))
}