package request

import (
	"encoding/gob"
	"fmt"
	"io"
	"github.com/materials-commons/materials/transfer"
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

// A IdentityMarshaler saves the data passed and returns it.
// It can be set to return an error instead. This is useful
// for testing.
type RequestMarshaler struct {
	request transfer.Request
	err  error
}

// NewIdentityMarshaler returns a new IdentityMarshaler
func NewRequestMarshaler() *RequestMarshaler {
	return &RequestMarshaler{}
}

// Marshal saves the data to be returned by the Unmarshal. If
// SetError has been called it instead returns the error passed
// to SetError and doesn't save the data.
func (m *RequestMarshaler) Marshal(data interface{}) error {
	if m.err != nil {
		return m.err
	}

	switch t := data.(type) {
	case *transfer.Request:
		m.request = *t
	default:
		return fmt.Errorf("Not a transfer.Request")
	}

	return nil
}

// Unmarshal returns the last data successfully passed to Marshal. If
// SetError has been called it instead returns the error passed to
// SetError and doesn't set the data.
func (m *RequestMarshaler) Unmarshal(data interface{}) error {
	if m.err != nil {
		return m.err
	}

	switch t := data.(type) {
	case *transfer.Request:
		*t = m.request
	default:
		fmt.Errorf("Not a transfer.Request")
	}

	return nil
}

// SetError sets the error that Marshal and Unmarshal should return.
func (m *RequestMarshaler) SetError(err error) {
	m.err = err
}

// ClearError clears the error so that Marshal and Unmarshal will no
// longer return an error when called.
func (m *RequestMarshaler) ClearError() {
	m.err = nil
}

// SetData will explicitly set the data rather than using Marshal. Useful
// in some test cases.
func (m *RequestMarshaler) SetData(data interface{}) {
	switch t := data.(type) {
	case *transfer.Request:
		m.request = *t
	}
}
