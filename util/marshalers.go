package util

import (
	"encoding/gob"
	"fmt"
	"github.com/materials-commons/mcfs/protocol"
	"io"
	"bytes"
)

var _ = fmt.Println

type ChannelReadWriter struct {
	c   chan []byte
	err error
}

func NewChannelReadWriter() *ChannelReadWriter {
	return &ChannelReadWriter{
		c: make(chan []byte),
	}
}

func (this *ChannelReadWriter) Write(bytes []byte) (n int, err error) {
	if this.err != nil {
		return 0, err
	}

	this.c <- bytes
	return len(bytes), nil
}

func (this *ChannelReadWriter) Read(bytes []byte) (n int, err error) {
	if this.err != nil {
		return 0, err
	}

	bytes = <-this.c
	return len(bytes), nil
}

// A GobMarshaler marshals and unmarshals data using Gob.
type GobMarshaler struct {
	*gob.Encoder
	*gob.Decoder
}

// NewGobMarshaler returns a new GobMarshaler.
func NewGobMarshaler(rw io.ReadWriter) *GobMarshaler {
	return &GobMarshaler{
		Encoder: gob.NewEncoder(rw),
		Decoder: gob.NewDecoder(rw),
	}
}

// Marshal marshals the data using gob.Encode.
func (m *GobMarshaler) Marshal(data interface{}) error {
	return m.Encode(data)
}

// Unmarshal unmarshals the data using gob.Decode.
func (m *GobMarshaler) Unmarshal(data interface{}) error {
	return m.Decode(data)
}

/* ******************************************************************* */

type ChannelMarshaler struct {
	request  chan protocol.Request
	response chan protocol.Response
	err      error
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewChannelMarshaler() *ChannelMarshaler {
	var buf bytes.Buffer
	return &ChannelMarshaler{
		request:  make(chan protocol.Request),
		response: make(chan protocol.Response),
		encoder: gob.NewEncoder(&buf),
		decoder: gob.NewDecoder(&buf),
	}
}

func (m *ChannelMarshaler) Marshal(data interface{}) error {
	if m.err != nil {
		return m.err
	}
	switch t := data.(type) {
	case *protocol.Request:
		m.request <- *t
	case protocol.Request:
		m.request <- t
	case *protocol.Response:
		m.response <- *t
	case protocol.Response:
		m.response <- t
	}
	return nil
}

func (m *ChannelMarshaler) Unmarshal(data interface{}) error {
	if m.err != nil {
		return m.err
	}

	select {
	case req := <-m.request:
		switch t := data.(type) {
		case *protocol.Request:
			*t = req
		default:
			return fmt.Errorf("Request data needed")
		}
	case resp := <-m.response:
		switch t := data.(type) {
		case *protocol.Response:
			*t = resp
		default:
			return fmt.Errorf("Response data needed")
		}
	}
	return nil
}

func (m *ChannelMarshaler) SetError(err error) {
	m.err = err
}

func (m *ChannelMarshaler) ClearError() {
	m.err = nil
}

/* ******************************************************************* */

// A IdentityMarshaler saves the data passed and returns it.
// It can be set to return an error instead. This is useful
// for testing.

type RequestResponseMarshaler struct {
	request  protocol.Request
	response protocol.Response
	err      error
}

// NewIdentityMarshaler returns a new IdentityMarshaler
func NewRequestResponseMarshaler() *RequestResponseMarshaler {
	return &RequestResponseMarshaler{}
}

// Marshal saves the data to be returned by the Unmarshal. If
// SetError has been called it instead returns the error passed
// to SetError and doesn't save the data.
func (m *RequestResponseMarshaler) Marshal(data interface{}) error {
	if m.err != nil {
		return m.err
	}

	switch t := data.(type) {
	case *protocol.Request:
		m.request = *t
	case *protocol.Response:
		m.response = *t
	default:
		return fmt.Errorf("Not a valid type")
	}

	return nil
}

// Unmarshal returns the last data successfully passed to Marshal. If
// SetError has been called it instead returns the error passed to
// SetError and doesn't set the data.
func (m *RequestResponseMarshaler) Unmarshal(data interface{}) error {
	if m.err != nil {
		return m.err
	}

	switch t := data.(type) {
	case *protocol.Request:
		*t = m.request
	case *protocol.Response:
		*t = m.response
	default:
		return fmt.Errorf("Not a valid type")
	}

	return nil
}

// SetError sets the error that Marshal and Unmarshal should return.
func (m *RequestResponseMarshaler) SetError(err error) {
	m.err = err
}

// ClearError clears the error so that Marshal and Unmarshal will no
// longer return an error when called.
func (m *RequestResponseMarshaler) ClearError() {
	m.err = nil
}

// SetData will explicitly set the data rather than using Marshal. Useful
// in some test cases.
func (m *RequestResponseMarshaler) SetData(data interface{}) {
	switch t := data.(type) {
	case *protocol.Request:
		m.request = *t
	case *protocol.Response:
		m.response = *t
	}
}
