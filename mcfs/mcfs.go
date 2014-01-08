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
	"github.com/jessevdk/go-flags"
	"github.com/materials-commons/materials/mcfs/request"
	_ "github.com/materials-commons/materials/transfer"
	"net"
	"os"
)

type ServerOptions struct {
	Port     uint   `long:"server-port" description:"The port the server listens on" default:"35862"`
	Bind     string `long:"bind" description:"Address of local interface to listen on" default:"localhost"`
	PrintPid bool   `long:"print-pid" description:"Prints the server pid to stdout"`
}

type DatabaseOptions struct {
	Connection string `long:"db-connect" description:"The host/port to connect to database on" default:"localhost:28015"`
	Name       string `long:"db" description:"Database to use" default:"materialscommons"`
}

type Options struct {
	Server   ServerOptions `group:"Server Options"`
	Database DatabaseOptions
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	listener, err := createListener(opts.Server.Bind, opts.Server.Port)
	if err != nil {
		os.Exit(1)
	}

	if opts.Server.PrintPid {
		fmt.Println(os.Getpid())
	}

	acceptConnections(listener, opts.Database.Connection, opts.Database.Name)
}

func createListener(host string, port uint) (*net.TCPListener, error) {
	service := fmt.Sprintf("%s:%d", host, port)
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

func acceptConnections(listener *net.TCPListener, dbAddress, dbName string) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		session, _ := r.Connect(map[string]interface{}{
			"address":  dbAddress,
			"database": dbName,
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
