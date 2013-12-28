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
	"github.com/materials-commons/materials/model"
	"github.com/materials-commons/materials/transfer"
	"net"
	"os"
	"path/filepath"
	"strings"
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
	conn net.Conn
	db   *rethink.DB
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	command := getCommand(conn)
	if !transfer.ValidType(command.Type) {
		fmt.Println("Invalid command:", command.Type)
		return
	}

	handler, err := createHandler(command, conn)
	switch {
	case err != nil:
		fmt.Println("Error creating connection handler", err)
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
		conn:    conn,
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
	response := transfer.SendStartResponse{
		Offset:     0,
		DataFileID: h.DataFile.ID,
	}
	datafile, err := model.GetDataFile(h.DataFile.ID, h.db.Session)

	if err != nil || !ownerGaveAccessTo(datafile.Owner, h.Header.User, h.db.Session) {
		return
	}

	if h.DataFile.Checksum == datafile.Checksum && h.DataFile.Size != datafile.Size {
		response.Offset = datafile.Size
	}

	encoder := gob.NewEncoder(h.conn)
	encoder.Encode(&response)
	decoder := gob.NewDecoder(h.conn)
	/*
	* Need to open file and write to it.
	 */
	filename := createPath(h.DataFile.ID)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	for {
		buf := &transfer.FileBlock{}
		decoder.Decode(buf)
		if n, err := f.Write(buf.Bytes); err != nil {
			// Do something, check n vs bytes
			fmt.Println(n)
		}
		if buf.Done {
			break
		}
	}

	/*
	* Set the size in rethinkdb to the number of bytes we have written
	* plus the size of the file that is already on the isilon.
	 */
}

func createPath(datafileId string) string {
	pieces := strings.Split(datafileId, "-")
	dirpath := filepath.Join("/mcfs/data/materialscommons", pieces[1][0:2], pieces[1][2:4])
	os.MkdirAll(dirpath, 0600)
	return filepath.Join(dirpath, datafileId)
}

// hasAccess checks to see if the user making the request has access to the
// particular datafile. Access is determined as follows:
// 1. if the user and the owner of the file are the same return true (has access).
// 2. Get a list of all the users groups for the file owner.
//    For each user in the user group see if teh requesting user
//    is included. If so then return true (has access).
// 3. None of the above matched - return false (no access)
func ownerGaveAccessTo(owner, user string, session *r.Session) bool {
	// Check if user and file owner are the same
	if user == owner {
		return true
	}

	// Get the file owners usergroups
	rql := r.Table("usergroups").Filter(r.Row.Field("owner").Eq(owner))
	groups, err := model.MatchingUserGroups(rql, session)
	if err != nil {
		return false
	}

	// For each usergroup go through its list of users
	// and see if they match the requesting user
	for _, group := range groups {
		users := group.Users
		for _, u := range users {
			if u == user {
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
