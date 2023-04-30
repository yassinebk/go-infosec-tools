package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type router struct {
}

type logger struct {
	Inner http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start")

	l.Inner.ServeHTTP(w, r)

	log.Println("end")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World ")
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/a":
		fmt.Fprintf(w, "Executing a")
	case "/b":
		fmt.Fprintf(w, "Executing b")
	case "/c":
		fmt.Fprintf(w, "Executing c")
	default:
		http.Error(w, "404 not found", 404)
	}

}

func main() {
	f := http.HandlerFunc(hello)

	r := mux.NewRouter()

	r.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Moew meow foo")
	}).Methods("GET").Host("www.foo.com")

	r.HandleFunc("/users/{user}", func(w http.ResponseWriter, req *http.Request) {
		user := mux.Vars(req)["user"]
		fmt.Fprintf(w, "hi %s\n", user)
	}).Methods("GET")

	l := logger{Inner: f}

	http.ListenAndServe(":8000", &l)

}
