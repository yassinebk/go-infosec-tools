package main

/**
	the goal here is to write an echo server that writes and read data from a socket
Steps:
	1. Establish a tcp connection
	3. Accepting the connection
	2. Attach my reader and writer
*/

import (
	"bufio"
	"log"
	"net"
)

func echo(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)

		s, err := reader.ReadString('\n')

		if err != nil {

			log.Fatalln("Unable to read data")
		}

		log.Println("Writing data to socket")

		// Hooking the writer to the socket
		writer := bufio.NewWriter(conn)

		if _, err := writer.WriteString(s); err != nil {

			log.Fatalln("Unable to write data")
		}

		writer.Flush()
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
