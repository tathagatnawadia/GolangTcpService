package main

import (
	"log"
	"fmt"
)

import (
	"./Entities"
	"./RelayServer"
)

func main() {

	// Create the hub network and attach it to the server
	unityNetwork, err := Entities.NewNetwork()
	if err != nil {
		log.Fatal(err)
	}

	server, err := RelayServer.NewServer("tcp", ":6666", unityNetwork)
	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	// Start the server
    server.Run()

	// Detect closed clients
	go server.DetectDisconnectedClients()

	// Handle the clients
	fmt.Println("Starting the application at localhost:6666")
	fmt.Println("Join the hub via \n$nc localhost 6666")

	server.HandleConnections()
}