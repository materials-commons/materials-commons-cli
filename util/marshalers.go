package util

import (
	"encoding/gob"
	"fmt"
	"github.com/materials-commons/materials/transfer"
	"io"
)

var _ = fmt.Println

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
	request  chan transfer.Request
	response chan transfer.Response
	err      error
}

func NewChannelMarshaler() *ChannelMarshaler {
	return &ChannelMarshaler{
		request:  make(chan transfer.Request),
		response: make(chan transfer.Response),
	}
}

func (m *ChannelMarshaler) Marshal(data interface{}) error {
	if m.err != nil {
		return m.err
	}
	switch t := data.(type) {
	case *transfer.Request:
		m.request <- *t
	case transfer.Request:
		m.request <- t
	case *transfer.Response:
		m.response <- *t
	case transfer.Response:
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
		case *transfer.Request:
			*t = req
		default:
			return fmt.Errorf("Request data needed")
		}
	case resp := <-m.response:
		switch t := data.(type) {
		case *transfer.Response:
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
	request  transfer.Request
	response transfer.Response
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
	case *transfer.Request:
		m.request = *t
	case *transfer.Response:
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
	case *transfer.Request:
		*t = m.request
	case *transfer.Response:
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
	case *transfer.Request:
		m.request = *t
	case *transfer.Response:
		m.response = *t
	}
}
