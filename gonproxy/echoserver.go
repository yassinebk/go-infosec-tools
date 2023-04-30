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

func echo(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 512)

	for {
		size, err := conn.Read(b[0:])

		if err == io.EOF {

			log.Fatalln("Connection closed by client")
			break
		}
		if err != nil {
			log.Fatalln("Unable to read data")
			break
		}

		log.Printf("Read %d bytes from socket\n", size)

		log.Println("Writing data to socket")

		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("Unable to write data")
			break
		}
		log.Printf("Written %d bytes to socket\n", size)
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
