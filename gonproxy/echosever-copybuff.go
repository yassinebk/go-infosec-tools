package main

/**
	the goal here is to write an echo server that writes and read data from a socket
Steps:
	1. Establish a tcp connection
	3. Accepting the connection
	2. Attach my reader and writer
*/

import (
	"io"
	"log"
	"net"
)

// The simplest one heeere
func echo(conn net.Conn) {
	defer conn.Close()

	for {
		if _, err := io.Copy(conn, conn); err != nil {
			log.Fatalln("Unable to copy data")
		}
	}

}
func main() {

	listener, err := net.Listen("tcp", ":4000")

	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	log.Println("Listening on 0.0.0.0:4000")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		log.Println("Accepted connection from client")

		go echo(conn)
	}

}
