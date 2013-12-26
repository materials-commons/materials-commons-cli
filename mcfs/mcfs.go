/*
 * This package implements the Materials Commons File Server service. This
 * service provides upload/download of datafiles from the Materials Commons
 * repository.

 * The protocol for file uploads looks as follows:
 *     1. The client sends the size, checksum and path. If the file
 *        is an existing file then it also sends the DataFileID for
 *        the file.
 *
 *     2. If the server receives a DataFileID it checks the size
 *        and checksum against what was sent. If the checksums
 *        match and the sizes are different then its a partially
 *        completed upload. If the checksums are different then
 *        its a new upload.
 *
 *     3. The server sends back the DataFileID. It will create a
 *        new DataFileID or send back an existing depending on
 *        whether its a new upload or an existing one.
 *
 *     4. The server will tell the client the offset to start
 *        sending data from. For a new upload this will be at
 *        position 0. For an existing one it will be the offset
 *        to restart the upload.
 *
 * The protocol for file downloads looks as follows:
 *
 */
package main

import (
	"encoding/gob"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/transfer"
	"net"
	"os"
)

const Port = "35862"

func main() {
	service := "0.0.0.0:" + Port
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println("Resolve error:", err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Listen error:", err)
		os.Exit(1)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

type FileTransferHeader2 struct {
	Size  int
	Bytes []byte
}

type commandHandler struct {
	*transfer.Command
	net.Conn
	session *r.Session
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	command := &transfer.Command{}
	decoder.Decode(command)

	session, err := r.Connect(map[string]interface{}{
		"address":  "localhost:30815",
		"database": "materialscommons",
	})

	if err != nil {
		fmt.Println("Unable to connect to database")
		return
	}

	handler := &commandHandler{
		Command: command,
		Conn:    conn,
		session: session,
	}

	if handler.validApiKey() {
		switch command.Type {
		case transfer.Upload:
			handler.upload()
		case transfer.Download:
			handler.download()
		case transfer.Move:
			handler.move()
		case transfer.Delete:
			handler.delete()
		default:
			fmt.Println("Unknown command type: %d", command.Type)
		}
	}
	/*
		fth := &FileTransferHeader2{}
		decoder.Decode(fth)
		fmt.Printf("Received %#v\n", fth)
		fmt.Println(string(fth.Bytes))
	*/
}

func (h *commandHandler) validApiKey() bool {
	apikey, err := h.queryUser()
	switch {
	case err != nil:
		return false
	case apikey != h.Header.ApiKey:
		return false
	default:
		return true
	}
}

func (h *commandHandler) queryUser() (string, error) {
	result, err := r.Table("users").Get(h.Header.Owner).RunRow(h.session)
	if err != nil || result.IsNil() {
		return "", fmt.Errorf("Unknown user '%s'", h.Header.Owner)
	}

	var response map[string]interface{}
	result.Scan(&response)
	apikey := response["apikey"].(string)
	return apikey, nil
}

func (h *commandHandler) upload() {

}

func (h *commandHandler) download() {

}

func (h *commandHandler) move() {

}

func (h *commandHandler) delete() {

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
	b := "Hello world Bytes"
	encoder.Encode([]byte(b))
}
