package main

import (
	"log"
	"net"
)



func main() {
	listener, err := net.Listen("tcp", ":4000")

	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	for {
		conn, err := listener.Accept()

		log.Println("Accepted connection from client")

		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		go Handle(conn)
	}
}
