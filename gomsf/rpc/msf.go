package rpc

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/vmihailenco/msgpack.v2"
)

/*
*

	The data is passed as an array, not a map, so rather than expecting data in key/value format, the RPC interface expects the data as a positional array of values.
*/
type sessionListReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Token    string
}

type SessionListRes struct {
	ID          uint32 `msgpack:",omitempty"`
	Type        string `msgpack:"type"`
	TunnelLocal string `msgpack:"tunnel_local"`
	TunnelPeer  string `msgpack:"tunnel_peer"`
	ViaExploit  string `msgpack:"via_exploit"`
	ViaPayload  string `msgpack:"via_payload"`
	Description string `msgpack:"desc"`
	Info        string `msgpack:"info"`
	Workspace   string `msgpack:"workspace"`
	SessionHost string `msgpack:"session_host"`
	SessionPort int    `msgpack:"session_port"`
	Username    string `msgpack:"username"`
	UUID        string `msgpack:"uuid"`
	ExploitUUID string `msgpack:"exploit_uuid"`
}

type loginReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Username string
	Password string
}

type loginRes struct {
	Result       string `msgpack:"result"`
	Token        string `msgpack:"token"`
	Error        bool   `msgpack:"error"`
	ErrorClass   string `msgpack:"error_class"`
	ErrorMEssage string `msgpack:"error_message"`
}

type logoutReq struct {
	_msgpack    struct{} `msgpack:",asArray"`
	Method      string
	Token       string
	LogoutToken string
}

type LogoutRes struct {
	Result string `msgpack:"result"`
}

type Metasploit struct {
	host  string
	user  string
	pass  string
	token string
}

func New(host, user, pass string) (*Metasploit, error) {
	msf := &Metasploit{
		host: host,
		user: user,
		pass: pass,
	}

	err := msf.Login()
	if err != nil {
		return nil, err
	}

	return msf, nil
}

// A simple wrapper on the msg pack rpc call
func (msf *Metasploit) send(req interface{}, res interface{}) error {
	buf := new(bytes.Buffer)
	msgpack.NewEncoder(buf).Encode(req)

	dest := fmt.Sprintf("http://%s/api", msf.host)


	r, err := http.Post(dest, "binary/message-pack", buf)

	if err != nil {

		return err
	}

	defer r.Body.Close()

	if err := msgpack.NewDecoder(r.Body).Decode(res); err != nil {
		return err
	}

	return nil

}

func (msf *Metasploit) Login() error {

	ctx := &loginReq{
		Method:   "auth.login",
		Username: msf.user,
		Password: msf.pass,
	}

	res := &loginRes{}

	err := msf.send(ctx, res)

	if err != nil {
		log.Panicln("Error logging in Metasploit: ", err)

		return err
	}

	msf.token = res.Token

	return nil

}

func (msf *Metasploit) Logout() error {

	ctx := &logoutReq{
		Method:      "auth.logout",
		Token:       msf.token,
		LogoutToken: msf.token,
	}

	res := &LogoutRes{}

	err := msf.send(ctx, res)

	if err != nil {
		log.Panicln("Error logging out from Metasploit: ", err)

		return err
	}

	msf.token = ""

	return nil

}

func (msf *Metasploit) SessionList() (map[uint32]SessionListRes, error) {

	ctx := &sessionListReq{
		Method: "session.list",
		Token:  msf.token,
	}

	res := make(map[uint32]SessionListRes)

	err := msf.send(ctx, &res)

	if err != nil {
		log.Panicln("Error Listing sessions from Metasploit: ", err)

		return nil, err
	}

	for id, session := range res {
		session.ID = id
		res[id] = session
	}
	return res, nil

}


