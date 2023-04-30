package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type FooReader struct{}

type FooWriter struct{}

func (f FooReader) Read(p []byte) (n int, err error) {

	fmt.Printf("in < ")
	return os.Stdin.Read(p)
}

func (f FooWriter) Write(p []byte) (n int, err error) {

	fmt.Printf("out > ")
	return os.Stdout.Write(p)
}

func main() {

	var (
		reader FooReader
		writer FooWriter
	)

	// Copy reads from reader and writes to writer.
	if _, err := io.Copy(&writer, &reader); err != nil {
		log.Fatalln("Unable to read/write data")
	}
}
