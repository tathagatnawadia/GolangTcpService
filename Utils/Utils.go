package Utils

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

func SendResponse(res string, conn io.Writer) {
	fmt.Printf("[%s][RESPONSE] %s\n", time.Now().Format(time.RFC3339), res)
	conn.Write([]byte(res + "\n"))
}

func SendPrompt(res string, conn io.Writer) {
	conn.Write([]byte(res))
}

func SendBroadcast(res string, from int, conn io.Writer) {
	res = "\nMessage from " + strconv.Itoa(from) + " : " + res + "\n"
	conn.Write([]byte(res))
}

func PrintHelpText(conn io.Writer) {
	res := "----------HELP---------\n"
	res += ">> LIST\n"
	res += ">> IDENTIFY\n"
	res += ">> RELAY #How you doing user 1,3 and 4 ?? #1,3,4\n"
	res += ">> EXIT\n"
	res += "-----------------------\n"
	conn.Write([]byte(res))
}
