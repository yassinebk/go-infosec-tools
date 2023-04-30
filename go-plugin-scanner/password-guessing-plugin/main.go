package main

import (
	"fmt"
	scanner "go-plugin-scanner/scanner"
	"log"
	"net/http"
)

var Users = []string{"admin", "manager", "tomcat"}
var Passwords = []string{"admin", "manager", "tomcat", "password"}

type TomcatChecker struct {
}

func (c *TomcatChecker) Check(host string, port uint64) *scanner.Results {
	var (
		resp   *http.Response
		err    error
		res    *scanner.Results
		client *http.Client
		req    *http.Request
	)

	log.Println("Checking for tomcat Manager creds ...")

	res = new(scanner.Result)
	url := fmt.Sprintf("http://%s:%d/manager/html", host, port)
	if resp, err = http.Head(url); err != nil {
		log.Println(err)
		return res
	}

	log.Println("Host responded to /manager/html request")

	if resp.StatusCode != http.StatusUnauthorized || resp.Header.Get("WWW-Authenticate") == "" {
		log.Println("Host did not respond with 401, it doesn't require basic authentication")
		return res
	}

	log.Println("Host requires basic authentication")
	client = new(http.Client)

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Println(err)
		return res
	}

	for _, user := range Users {
		for _, password := range Passwords {
			req.SetBasicAuth(user, password)
			if resp, err = client.Do(req); err != nil {
				log.Println(err)
				continue
			}

			if resp.StatusCode == http.StatusOK {
				res.Vulnerable = true
				res.Details = fmt.Sprintf("Valid credentials found - %s:%s", user, password)
				return res
			}
		}

	}
	return res

}

func New() scanner.Checker {
	return new(TomcatChecker)
}
