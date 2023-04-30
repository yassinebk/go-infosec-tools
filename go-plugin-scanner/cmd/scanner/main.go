package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"plugin"

	scanner "go-plugin-scanner/cmd"
)

const PluginsDir = "../../plugins/"

func main() {
	var (
		files  []os.FileInfo
		err   error
		p     *plugin.Plugin
		n     plugin.Symbol
		check scanner.Checker
		res   *scanner.Result
	)

	if files, err = ioutil.ReadDir(PluginsDir); err != nil {
		panic(err)
	}

	for idx := range files {
		fmt.Println("Loading plugin: ", files[idx].Name())
		if p, err = plugin.Open(PluginsDir + files[idx].Name()); err != nil {
			panic(err)
		}

		if n, err = p.Lookup("New"); err != nil {
			panic(err)
		}

		newFunc, ok := n.(func() scanner.Checker)
		if !ok {
			panic("Invalid plugin. New() function not found")
		}

		check = newFunc()

		res = check.Check(os.Args[1], os.Args[2])
		if res.Vulnerable {
			fmt.Printf("Vulnerable: %s:", res.Details)
		} else {
			fmt.Println("Not vulnerable")
		}

	}
}
