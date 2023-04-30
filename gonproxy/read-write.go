package main

import (
	"fmt"
	"log"
	"os"
)

type FooReader struct{}

type FooWriter struct{}

func (f FooReader) Read(p []byte) (n int, err error) {

	fmt.Println("in <")
	return os.Stdin.Read(p)
}

func (f FooWriter) Write(p []byte) (n int, err error) {

	fmt.Println("out >")
	return os.Stdout.Write(p)
}

func main() {

	var (
		reader FooReader
		writer FooWriter
	)
	input := make([]byte, 4096)

	s, err := reader.Read(input)

	if err != nil {
		log.Fatalln("Unable to read data")
	}

	fmt.Printf("Read %d bytes from stdin", s)

	s, err = writer.Write(input)

	if err != nil {
		log.Fatalln("Unable to write data")
	}

	fmt.Printf("Written %d bytes from stdout]\n", s)

}
