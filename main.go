package main

import (
	"github.com/joho/godotenv"
	"net"
	"os"
	"log"
	"strings"
	"bufio"
)

import (
	"./Utils"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	Utils.SendPrompt("COMMAND : ", conn)
	scanNodeData := bufio.NewScanner(conn)

	for {
		for scanNodeData.Scan() {
			params := strings.Split(scanNodeData.Text(), " ")
			command := params[0]

			switch strings.Trim(strings.ToLower(command), " ") {
		        case "identify":
		            Utils.SendResponse("You asked to be identified", conn)
		        case "list":
		            Utils.SendResponse("You asked to list all active", conn)
		        case "relay":
		        	Utils.SendResponse("You asked to relay your message to other users", conn)
		        default:
		        	Utils.SendResponse("UNKWN Command : "+command, conn)
		    }
			Utils.SendPrompt("COMMAND : ", conn)
		}
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	/*
	------------------------------------------------------
	Purpose : Open a tcp connection for nodes to join via nc localhost 6666
	------------------------------------------------------
	*/
	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	/*
	------------------------------------------------------
	Purpose : Handle indivisual tcp connections and also close dead tcp connections
	------------------------------------------------------
	*/

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}