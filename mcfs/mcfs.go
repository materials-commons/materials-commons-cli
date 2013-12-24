package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

const Port = ":35862"

func main() {
	service := "0.0.0.0" + Port
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println("Resolve error:", err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Listen error:", err)
	}
	go client()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		handleConnection(conn)
		os.Exit(0)
	}
}

type FileTransferHeader2 struct {
	Size int
	Bytes []byte
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	fth := &FileTransferHeader2{}
	decoder.Decode(fth)
	fmt.Printf("Received %#v\n", fth)
	fmt.Println(string(fth.Bytes))
}

func client() {
	conn, err := net.Dial("tcp", "localhost"+Port)
	if err != nil {
		fmt.Println("Error on client connect", err)
		os.Exit(1)
	}
	encoder := gob.NewEncoder(conn)
	fth := &FileTransferHeader2{1, []byte("Hello World")}
	encoder.Encode(fth)
}
