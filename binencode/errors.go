package binencode

import (
	"errors"
)

// ErrOutputSize is returned if the output buffer capacity is too small.
var ErrOutputSize = errors.New("binencode: output buffer capacity too small")

// ErrInputSize is returned if the input buffer length is too small.
var ErrInputSize = errors.New("binencode: input buffer length too small")

// ErrSliceSize is returned if the output slice is too small.
var ErrSliceSize = errors.New("binencode: output slice too small")

// ErrSliceExpected is returned if the slice has an unexpected length.
var ErrSliceExpected = errors.New("binencode: slice has unexpected length")

// ErrSliceExpectedLong is returned if the slice has an unexpected long length.
var ErrSliceExpectedLong = errors.New("binencode: slice has unexpected long length")

// ErrType is returned if an unexpected type is encountered.
var ErrType = errors.New("binencode: unexpected type encountered")

// ErrNil is returned when writing a type to a nil value.
var ErrNil = errors.New("binencode: cannot write type to nil value")
