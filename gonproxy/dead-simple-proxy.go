package main

import (
	"io"
	"log"
	"net"
)

func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		log.Fatalln("Unable to connect to remote server")
	}

	defer dst.Close()

	// Copy the request from the client to the remote server
	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln("Unable to copy data")
		}
	}()

	// Copy the response from the remote server back to the client
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
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

		go handle(conn)
	}
}
