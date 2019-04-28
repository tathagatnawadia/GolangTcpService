package main

import (
	"fmt"
	"log"
	"./Entities"
	"./RelayServer"
)

const (
	TCP_ADDR = "0.0.0.0:6666"
	TCP_PORT = "6666"
	TCP_HOST = "localhost"
)

func main() {

	// Create the hub network and attach it to the server
	myNetwork, err := Entities.NewNetwork()
	if err != nil {
		log.Fatal(err)
	}

	server, err := RelayServer.NewServer("tcp", TCP_ADDR, myNetwork)
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	// Start the server
	server.Run()

	// Detect closed clients
	go server.DetectDisconnectedClients()

	// Handle the clients
	fmt.Printf("Starting the application at %s\n", TCP_ADDR)
	fmt.Printf("Join the hub via \n$nc %s %s\n", TCP_HOST, TCP_PORT)

	server.HandleConnections()
}
