package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	_, err := net.Dial("tcp", "scanme.nmap.org:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connection successful!")

	var wg sync.WaitGroup
	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", fmt.Sprintf("scanme.nmap.org:%d", j))

			if err != nil {
				return
			}

			conn.Close()

			fmt.Println("port", j, "is open")

		}(i)

	}
	wg.Wait()

}
