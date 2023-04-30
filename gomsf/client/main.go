package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yassinebk/gomsf/rpc"
)

/*
Steps to build our HTTP Client that will communicated with Shodan API to gather information about the target.

	1. Review the service’s API documentation.
	2. Design a logical structure for the code in order to reduce complexity
	and repetition.
	3. Define request or response types, as necessary, in Go.
	4. Create helper functions and types to facilitate simple initialization,
	authentication, and communication to reduce verbose or repetitive
	logic.
	5. Build the client that interacts with the API consumer functions and
types.

*/

func main() {

	host := os.Getenv("MSF_HOST")
	pass := os.Getenv("MSF_PASS")

	user := "msf"

	if host == "" || pass == "" {
		log.Fatal("Please set the environment variables MSF_HOST and MSF_PASS")
	}

	msf, err := rpc.New(host, user, pass)
	if err != nil {
		log.Fatalln("Error creating the client", err)
	}

	defer msf.Logout()

	fmt.Println("Sessions 1")
	sessions, err := msf.SessionList()

	if err != nil {
		log.Fatalln("Error listing sessions", err)
	}

	for _, session := range sessions {
		fmt.Printf("%5d %s\n", session.ID, session.Info)
	}

}
