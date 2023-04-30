package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yassinebk/goshodattp/shodan"
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

	if len(os.Args) != 2 {

		log.Fatalln("Usage shodan <query>")
	}

	apiKey := os.Getenv("SHODAN_API_KEY")

	s := shodan.New(apiKey)

	info, err := s.APIInfo()

	if err != nil {
		log.Fatalln("Error api infos, please check your api key", err)
	}

	fmt.Printf("Query Credits: %d\n Scan Credits %d \n", info.QueryCredits, info.ScanCredits)

	hostSearch, err := s.HostSearch(os.Args[1])

	if err != nil {
		log.Panicln(err)
	}

	for _, host := range hostSearch.Matches {

		fmt.Printf("IP: %18s, PORT: %8d \n", host.IPstring, host.Port)
	}

}
