package mcc

import (
	"context"
	"os"
	"time"
)

type FileStatus struct {
	Path        string
	ProjectPath string
	Status      string
	FInfo       os.FileInfo
}

type FileStatusHandlerFn func(fs *FileStatus)

type FileStatusReceiver struct {
	in        chan *FileStatus
	HandlerFn FileStatusHandlerFn
}

func NewFileStatusReceiver(handler FileStatusHandlerFn) *FileStatusReceiver {
	return &FileStatusReceiver{
		in:        make(chan *FileStatus, 1),
		HandlerFn: handler,
	}
}

func (r *FileStatusReceiver) Run(c context.Context) {
	for {
		select {
		case fstatus := <-r.in:
			r.HandlerFn(fstatus)
		case <-c.Done():
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func (r *FileStatusReceiver) SendStatus(fs *FileStatus) {
	r.in <- fs
}
