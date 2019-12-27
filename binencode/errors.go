package binencode

import (
	"errors"
)

// ErrOutputSize is returned if the output buffer capacity is too small.
var ErrOutputSize = errors.New("binencode: output buffer capacity too small")

// ErrInputSize is returned if the input buffer length is too small.
var ErrInputSize = errors.New("binencode: input buffer length too small")

// ErrSlizeSize is returned if the output slize is too small.
var ErrSlizeSize = errors.New("binencode: output slize too small")

// ErrSlizeExpected is returned if the slice has an unexpected length.
var ErrSlizeExpected = errors.New("binencode: slice has unexpected length")

// ErrSlizeExpectedLong is returned if the slice has an unexpected long length.
var ErrSlizeExpectedLong = errors.New("binencode: slice has unexpected long length")

// ErrType is returned if an unexpected type is encountered.
var ErrType = errors.New("binencode: unexpected type encountered")

// ErrNil is returned when writing a type to a nil value.
var ErrNil = errors.New("binencode: cannot write type to nil value")
