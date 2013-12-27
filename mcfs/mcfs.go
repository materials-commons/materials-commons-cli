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
	"github.com/materials-commons/gohandy/rethink"
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
	db *rethink.DB
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	command := getCommand(conn)
	if !transfer.ValidType(command.Type) {
		return
	}

	handler, err := createHandler(command, conn)

	switch {
	case err != nil:
		fmt.Println("Unable to connect to database")
	default:
		handler.doCommand()
	}
}

func getCommand(conn net.Conn) *transfer.Command {
	command := &transfer.Command{}
	gob.NewDecoder(conn).Decode(command)
	return command
}

func createHandler(c *transfer.Command, conn net.Conn) (*commandHandler, error) {
	session, err := r.Connect(map[string]interface{}{
		"address":  "localhost:30815",
		"database": "materialscommons",
	})

	handler := &commandHandler{
		Command: c,
		Conn:    conn,
		db:      rethink.NewDB(session),
	}

	return handler, err
}

func (h *commandHandler) doCommand() {
	if !h.validApiKey() {
		fmt.Printf("Invalid apikey '%s' for user '%s'\n", h.Header.ApiKey, h.Header.User)
		return
	}

	switch h.Type {
	case transfer.Upload:
		h.upload()
	case transfer.Download:
		h.download()
	case transfer.Move:
		h.move()
	case transfer.Delete:
		h.delete()
	default:
		fmt.Println("Unknown command type: %d", h.Type)
	}
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
	result, err := h.db.Get("users", h.Header.User)
	if err != nil {
		return "", fmt.Errorf("Unknown user '%s'", h.Header.User)
	}

	apikey, found := result["apikey"]
	if found {
		return apikey.(string), nil
	} else {
		return "", nil
	}
}

func (h *commandHandler) upload() {
	if h.DataFile.ID != "" {
		h.uploadExisting()
	} else {
		h.uploadNew()
	}
}

func (h *commandHandler) uploadExisting() {
	datafile, err := h.db.Get("users", h.DataFile.ID)
	if err != nil || !h.hasAccess(datafile["owner"].(string)) {
		return
	}

}

// hasAccess checks to see if the user making the request has access to the
// particular datafile. Access is determined as follows:
// 1. if the user and the owner of the file are the same return true (has access).
// 2. Get a list of all the users groups for the file owner.
//    For each user in the user group see if teh requesting user
//    is included. If so then return true (has access).
// 3. None of the above matched - return false (no access)
func (h *commandHandler) hasAccess(owner string) bool {
	// Check if user and file owner are the same
	if h.Header.User == owner {
		return true
	}

	// Get the file owners usergroups
	rql := r.Table("usergroups").Filter(r.Row.Field("owner").Eq(owner))
	groups, err := h.db.GetAll(rql)
	if err != nil {
		return false
	}
	// For each usergroup go through its list of users
	// and see if they match the requesting user
	for _, group := range groups {
		users := group["users"].([]string)
		for _, user := range users {
			if h.Header.User == user {
				return true
			}
		}
	}
	return false
}

func (h *commandHandler) uploadNew() {

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
