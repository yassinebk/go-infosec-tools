package main

import (
	"fmt"
	"net"
)

func main() {
	_, err := net.Dial("tcp", "scanme.nmap.org:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connection successful!")

	for i := 1; i <= 1024; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("scanme.nmap.org:%d", i))

		if err != nil {
			continue
		}

		conn.Close()

		fmt.Println("port",i, "is open")

	}

}
