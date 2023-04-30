package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

// 1- Create a handler function to ingest an incoming query
// 2- Inspect the question in the query and extract the domain name
// 3- Identify the upstream DNS server correlating to the domain name
// 4- Exchange the question with the upstream DNS server and write the
// response to the client

func parse(filename string) (map[string]string, error) {
	records := make(map[string]string)
	fh, err := os.Open(filename)
	if err != nil {
		return records, err
	}

	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			return records, fmt.Errorf("%s is not valid line", line)
		}
		records[parts[0]] = parts[1]

	}
	return records, scanner.Err()
}

func main() {
	records, err := parse("proxy.config")
	if err != nil {
		panic(err)
	}

	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		fmt.Println("Incoming question %s", req.Question)
		if len(req.Question) < 1 {
			dns.HandleFailed(w, req)
			return
		}

		name := req.Question[0].Name
		fmt.Println("Name extracted from question %s", req.Question)

		parts := strings.Split(name, ".")
		if len(parts) > 1 {
			name = strings.Join(parts[len(parts)-2:], ".")
		}

		match, ok := records[name]

		if !ok {
			dns.HandleFailed(w, req)
			return
		}

		resp, err := dns.Exchange(req, match)
		if err != nil {
			dns.HandleFailed(w, req)
			return
		}

		if err := w.WriteMsg(resp); err != nil {
			dns.HandleFailed(w, req)
			return
		}
	})

	fmt.Printf("%+v\n", records)

	log.Fatalln(dns.ListenAndServe(":53", "udp", nil))
}
