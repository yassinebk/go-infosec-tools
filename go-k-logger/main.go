package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: (func(r *http.Request) bool { return true }),
	}

	listenAddr string
	wsAddr     string
	jsTemplate *template.Template
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":8080", "http listen address")
	flag.StringVar(&wsAddr, "ws-addr", "", "Address to connect to websocket")
	flag.Parse()

	var err error

	jsTemplate, err = template.ParseFiles("logger.js")

	if err != nil {
		log.Panicln(err)
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "", 500)
		return
	}

	defer conn.Close()

	fmt.Printf("Connection from %s\n", conn.RemoteAddr().String())

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			return
		}
		fmt.Printf("Message from %s: %s\n", conn.RemoteAddr().String(), string(msg))
	}

}

func serveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	jsTemplate.Execute(w, wsAddr)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/k.js", serveFile)
	r.HandleFunc("/ws", serveWS)
	log.Fatal(http.ListenAndServe(listenAddr, r))
}
