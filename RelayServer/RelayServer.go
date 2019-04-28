package RelayServer

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"../Entities"
	"../Utils"
)

type Server interface {
	Run() error
	Close() error
	HandleConnections() error
	HandleConnection(net.Conn)
	DetectDisconnectedClients()
}

type TCPServer struct {
	addr    string
	server  net.Listener
	network *Entities.Network
	lsofout string
	ticker  *time.Ticker
}

func (t *TCPServer) Run() (err error) {
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		return
	}
	return
}

func (t *TCPServer) Close() (err error) {
	return t.server.Close()
}

func (t *TCPServer) HandleConnections() (err error) {
	for {
		conn, err := t.server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go t.HandleConnection(conn)
		go t.removeDisconnectedClients(conn)
	}
}

func (t *TCPServer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	myClient := t.network.Register(conn)

	go myClient.ReceiveMessages()

	Utils.SendResponse("Welcome to the hub !", myClient.Handler)
	Utils.SendPrompt(">> ", myClient.Handler)
	scanNodeData := bufio.NewScanner(myClient.Handler)

	for {
		for scanNodeData.Scan() {
			myClient.AddToHistory(scanNodeData.Text())
			params := strings.Split(scanNodeData.Text(), "#")

			command := strings.Trim(strings.ToLower(params[0]), " ")
			fmt.Printf("[%s][REQUEST] %s\n", time.Now().Format(time.RFC3339), command)

			switch command {
			case "identify":

				user_id, ok := t.network.GetUserIdByConnection(myClient.Handler)
				if ok {
					Utils.SendResponse(strconv.Itoa(user_id), myClient.Handler)
				} else {
					panic(fmt.Sprintf("Not able to identify user_id in the network"))
				}

			case "list":
				Utils.SendResponse(t.network.GetActiveClients(myClient), myClient.Handler)

			case "relay":
				if len(params) != 3 {
					Utils.SendResponse("Type HELP for correct usage", myClient.Handler)
				} else {
					messageToBeSent := params[1]
					recievers := params[2]
					relayMessage := Entities.CreateRelayMessage(messageToBeSent, recievers, myClient.GetUserId())

					if relayMessage.ValidateMessageLength(1024) && relayMessage.ValidateRecieverCount(255) {
						t.network.SendRelayMessage(relayMessage, myClient)
						Utils.SendResponse("Message sent success", myClient.Handler)
					} else {
						Utils.SendResponse("Message not sent due to client voilations", myClient.Handler)
					}
				}

			case "exit":
				Utils.SendResponse("Bye", myClient.Handler)
				t.network.RemoveClientByConnection(myClient.Handler)

			default:
				Utils.PrintHelpText(myClient.Handler)
			}
			Utils.SendPrompt(">> ", myClient.Handler)
		}
	}
}

func (t *TCPServer) DetectDisconnectedClients() {
	for range t.ticker.C {
		lsof := exec.Command("lsof", "-p", strconv.Itoa(os.Getpid()), "-a", "-i", "tcp")
		grep := exec.Command("grep", "CLOSE_WAIT")
		lsofOut, _ := lsof.StdoutPipe()
		lsof.Start()
		grep.Stdin = lsofOut
		out, _ := grep.Output()
		t.lsofout = string(out)
	}
}

func (t *TCPServer) removeDisconnectedClients(conn net.Conn) {
	for range t.ticker.C {
		if strings.Index(t.lsofout, regexp.MustCompile(":[0-9]+$").FindAllStringSubmatch(conn.RemoteAddr().String(), -1)[0][0]+" ") > -1 {
			fmt.Println("Closing tcp connection - ", conn.RemoteAddr())
			t.network.RemoveClientByConnection(conn)
			return
		}
	}
}

func NewServer(protocol string, addr string, iNetwork *Entities.Network) (Server, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return &TCPServer{
			addr:    addr,
			network: iNetwork,
			lsofout: "",
			ticker:  time.NewTicker(time.Second * 5),
		}, nil
	}
	return nil, errors.New("Invalid protocol given")
}
