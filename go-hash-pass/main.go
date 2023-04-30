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
	if len(os.Args) != 5 {
		log.Fatalln("Usage: go run main.go <target/hosts> <user> <domain> <hash>")
	}

	buf, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatalln(err)
	}

	options := smb.Options{
		User:   os.Args[2],
		Domain: os.Args[3],
		Hash:   os.Args[4],
		Port:   445,
	}

	targets := bytes.Split(buf, []byte{'\n'})

	for _, target := range targets {
		options.Host = string(target)
		session, err := smb.NewSession(options, false)
		if err != nil {
			log.Printf("[-] Login failed %s: %s\n", options.Host, err)
			continue
		}

		defer session.Close()
		if session.IsAuthenticated {
			fmt.Println("[+] Login successful", options.Host)
		}
	}
}
