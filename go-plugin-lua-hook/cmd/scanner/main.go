package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	lua "github.com/yuin/gopher-lua"
)

const PluginsDir = "../../plugins/"

const LuaHttpTypeName = "http"

func main() {
	var (
		l     *lua.LState
		files []os.FileInfo
		err   error
		f     string
	)

	l = lua.NewState()
	defer l.Close()

	register(l)

	if files, err = ioutil.ReadDir(PluginsDir); err != nil {
		panic(err)
	}

	for idx := range files {
		fmt.Println("Loading plugin: ", files[idx].Name())
		f = fmt.Sprintf("%s/%s", PluginsDir, files[idx].Name())
		if err := l.DoFile(f); err != nil {

			log.Fatalln(err)
		}
	}
}

func head(l *lua.LState) int {

	var (
		host string
		port uint64
		path string
		resp *http.Response
		err  error
		url  string
	)

	host = l.CheckString(1)
	port = uint64(l.CheckInt64(2))
	path = l.CheckString(3)
	url = fmt.Sprintf("http://%s:%d/%s", host, port, path)

	if resp, err = http.Head(url); err != nil {
		l.Push(lua.LNumber(0))
		l.Push(lua.LBool(false))
		l.Push(lua.LString(fmt.Sprintf("Request failed: %s", err)))
		return 3
	}

	l.Push(lua.LNumber(resp.StatusCode))
	l.Push(lua.LBool(resp.Header.Get("WWW-Authenticate") != ""))
	l.Push(lua.LString(""))

	return 3
}

func get(l *lua.LState) int {
	var (
		host     string
		port     uint64
		path     string
		resp     *http.Response
		err      error
		url      string
		username string
		password string

		client *http.Client
		req    *http.Request
	)

	host = l.CheckString(1)
	port = uint64(l.CheckInt64(2))
	username = l.CheckString(4)
	password = l.CheckString(4)
	path = l.CheckString(5)
	url = fmt.Sprintf("http://%s:%d/%s", host, port, path)

	client = new(http.Client)

	if req, err = http.NewRequest("GET", url, nil); err != nil {
		l.Push(lua.LNumber(0))
		l.Push(lua.LBool(false))
		l.Push(lua.LString(fmt.Sprintf("GET: Unable to build request: %s", err)))
		return 3
	}

	if username != "" || password != "" {
		// Assume basic Auth is required since suer or password is set
		req.SetBasicAuth(username, password)
	}

	if resp, err = client.Do(req); err != nil {
		l.Push(lua.LNumber(0))
		l.Push(lua.LBool(false))
		l.Push(lua.LString(fmt.Sprintf("Request failed: %s", err)))
		return 3
	}

	l.Push(lua.LNumber(resp.StatusCode))
	l.Push(lua.LBool(false))
	l.Push(lua.LString(""))

	return 3

}

func register(l *lua.LState) {

	mt := l.NewTypeMetatable(LuaHttpTypeName)
	fmt.Println("Inside register")
	l.SetGlobal("http", mt)

	// static attributes
	l.SetField(mt, "head", l.NewFunction(head))
	l.SetField(mt, "get", l.NewFunction(get))

}
