package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

var (
	hostProxy = make(map[string]string)
	proxies   = make(map[string]*httputil.ReverseProxy)
)

func init() {
	hostProxy["attacker1.com"] = "http://localhost:8081"
	hostProxy["attacker2.com"] = "http://localhost:20080"
	for k, v := range hostProxy {
		remote, err := url.Parse(v)

		if err != nil {
			log.Fatal("Unable to parse proxy target")
		}

		proxies[k] = httputil.NewSingleHostReverseProxy(remote)
	}
}

func main() {
	r := mux.NewRouter()

	for host, proxy := range proxies {
		// Matching the proxy with host
		r.Host(host).Handler(proxy)
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}
