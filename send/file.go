package file

import ()

type FileMeta struct {
	Checksum []byte
	Size     int
	Name     string
	Id       string
	Owner    string
}
