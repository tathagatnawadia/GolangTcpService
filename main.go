package main

import (
	"github.com/joho/godotenv"
	"net"
	"os"
	"log"
	"strings"
	"bufio"
	"fmt"
	"strconv"
	"sync"
)

import (
	"./Entities"
	"./Utils"
)

var network = Entities.Network{make(map[int]*Entities.Client), make(map[net.Conn]int), 0, sync.Mutex{}}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	myClient := network.AddNewClient(conn)

	go myClient.ReceiveMessages()

	Utils.SendResponse("Welcome to the hub !", myClient.Handler)
	Utils.SendPrompt("ENTER COMMAND : ", myClient.Handler)
	scanNodeData := bufio.NewScanner(myClient.Handler)

	for {
		for scanNodeData.Scan() {
			myClient.AddToHistory(scanNodeData.Text())
			params := strings.Split(scanNodeData.Text(), "#")
			command := params[0]

			switch strings.Trim(strings.ToLower(command), " ") {
		        case "identify":

		        	user_id, ok := network.GetUserIdByConnection(myClient.Handler)
		        	if ok {
		        		Utils.SendResponse(strconv.Itoa(user_id), myClient.Handler)
		        	} else {
		        		panic(fmt.Sprintf("Not able to identify user_id in the network"))
		        	}
		            
		        case "list":
		        	Utils.SendResponse(network.GetActiveClients(myClient), myClient.Handler)

		        case "relay":
		        	messageToBeSent := params[1]
		        	recievers := params[2]
		        	relayMessage := Entities.CreateRelayMessage(messageToBeSent, recievers, myClient.User_id)

		        	if relayMessage.ValidateMessageLength(1024) && relayMessage.ValidateRecieverCount(255) {
		        		network.SendRelayMessage(relayMessage, myClient)
		        		Utils.SendResponse("Message sent success", myClient.Handler)
		        	} else {
		        		Utils.SendResponse("Message not sent due to client voilations", myClient.Handler)
		        	}

		        default:
		        	Utils.PrintHelpText(myClient.Handler)
		    }
			Utils.SendPrompt("ENTER COMMAND : ", myClient.Handler)
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

	fmt.Println("Starting the application at localhost:"+os.Getenv("ADDR"))
	fmt.Println("Join the hub via \n$nc localhost "+os.Getenv("ADDR"))

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}