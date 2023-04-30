package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
)

type badAuth struct {
	Username string
	Password string
}

func (b *badAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username != b.Username || password != b.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Creating context with username
	ctx := context.WithValue(r.Context(), "username", username)

	// Adding context to my request
	r = r.WithContext(ctx)

	next(w, r)

}

type trivial struct {
}

func (t *trivial) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	fmt.Println("Trivial middleware")
	next(w, r)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s", r.Context().Value("username"))
}

//  curl 'http://localhost:3000/hello?username=admin&password=admin
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	n := negroni.Classic()
	n.Use(&trivial{})
	n.Use(&badAuth{
		Username: "admin",
		Password: "admin",
	})
	n.UseHandler(r)
	http.ListenAndServe(":3000", n)
}
