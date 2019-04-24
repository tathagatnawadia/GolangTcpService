package main

import (
	"github.com/joho/godotenv"
	"net"
	"os"
	"log"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for{}
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