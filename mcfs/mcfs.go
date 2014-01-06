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
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/materials-commons/materials/mcfs/request"
	"github.com/materials-commons/materials/model"
	_ "github.com/materials-commons/materials/transfer"
	"net"
	"os"
)

const Port = "35862"

func main() {
	listener, err := createListener()
	if err != nil {
		os.Exit(1)
	}

	acceptConnections(listener)
}

func createListener() (*net.TCPListener, error) {
	service := "0.0.0.0:" + Port
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println("Resolve error:", err)
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Listen error:", err)
		return nil, err
	}

	return listener, nil
}

func acceptConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		session, _ := r.Connect(map[string]interface{}{
			"address":  "localhost:30815",
			"database": "materialscommons",
		})
		r := request.NewReqHandler(conn, session)
		go handleConnection(r, conn, session)
	}
}

func handleConnection(reqHandler *request.ReqHandler, conn net.Conn, session *r.Session) {
	defer conn.Close()
	defer session.Close()

	reqHandler.Run()
}

// ownerGaveAccessTo checks to see if the user making the request has access to the
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
