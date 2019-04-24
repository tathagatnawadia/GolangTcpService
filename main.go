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
	"time"
	"os/exec"
	"regexp"
)

import (
	"./Entities"
	"./Utils"
)

var network = Entities.Network{make(map[int]*Entities.Client), make(map[net.Conn]int), 0, sync.Mutex{}}
var lsofout = ""

var ticker = time.NewTicker(time.Second * 5)

func detectDisconnectedClients() {
	for range ticker.C {
		lsof := exec.Command("lsof", "-p", strconv.Itoa(os.Getpid()), "-a", "-i", "tcp")
		grep := exec.Command("grep", "CLOSE_WAIT")
		lsofOut, _ := lsof.StdoutPipe()
		lsof.Start()
		grep.Stdin = lsofOut
		out, _ := grep.Output()
		lsofout = string(out)
	}
}

func removeDisconnectedClients(conn net.Conn) {
	for range ticker.C {
		if strings.Index(lsofout, regexp.MustCompile(":[0-9]+$").FindAllStringSubmatch(conn.RemoteAddr().String(), -1)[0][0]+" ") > -1 {
			fmt.Println("Closing tcp connection - ", conn.RemoteAddr())
			network.RemoveClientByConnection(conn)
			return
		}
	}
}

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
	Purpose : Detect a closed connection and add it to the inventory
	------------------------------------------------------
	*/

	go detectDisconnectedClients()

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
		go removeDisconnectedClients(conn)
	}
}