// Package binencode provides functions to encode values into length-encoded slices. It can be used to work on
// secure memory. Only int16, int32, int64 and []byte are supported. Inserting and skipping zero bytes is supported.
package binencode

import (
	"encoding/binary"
	"errors"
	"reflect"
)

// Encode/decode into byteslices
// Type: 1 byte
// Length: 4 byte
// Data: [length]bytes (only for bytes)
// Encode/Decode move pointers into input/output slices further
// Encode/Decode check capacity of slices for ability
// - []byte (min,max)

var (
	ErrOutputSize        = errors.New("types: Output buffer capacity too small")
	ErrInputSize         = errors.New("types: Input buffer length too small")
	ErrSlizeSize         = errors.New("types: Output slize too small")
	ErrSlizeExpected     = errors.New("types: Slice has unexpected length")
	ErrSlizeExpectedLong = errors.New("types: Slice has unexpected long length")
	ErrType              = errors.New("types: Unexpected type encountered")
	ErrNil               = errors.New("types: Cannot write type to nil value")
)

const Encode16Size = 3

// EncodeInt16 encodes an int16 into b. It returns b, advanced by the data written, as well as the length of data written.
func EncodeInt16(i int16, b []byte) (output []byte, n int, err error) {
	if cap(b) < Encode16Size {
		return b, 0, ErrOutputSize
	}
	b = b[0:Encode16Size]
	b[0] = 0x01
	binary.BigEndian.PutUint16(b[1:Encode16Size], uint16(i))
	return b[Encode16Size:], Encode16Size, nil
}

// DecodeInt16 decodes an Int16 from b, it advances and returns b as well as the number of bytes read.
func DecodeInt16(b []byte, i *int16) (output []byte, n int, err error) {
	if len(b) < Encode16Size {
		return b, 0, ErrInputSize
	}
	if b[0] != 0x01 {
		return b, 0, ErrType
	}
	if i == nil {
		return b, 0, ErrNil
	}
	*i = int16(binary.BigEndian.Uint16(b[1:Encode16Size]))
	return b[Encode16Size:], Encode16Size, nil
}

const Encode32Size = 5

// EncodeInt32 encodes an int32 into b. It returns b, advanced by the data written, as well as the length of data written.
func EncodeInt32(i int32, b []byte) (output []byte, n int, err error) {
	if cap(b) < Encode32Size {
		return b, 0, ErrOutputSize
	}
	b = b[0:Encode32Size]
	b[0] = 0x02
	binary.BigEndian.PutUint32(b[1:Encode32Size], uint32(i))
	return b[Encode32Size:], Encode32Size, nil
}

// DecodeInt32 decodes an Int32 from b, it advances and returns b as well as the number of bytes read.
func DecodeInt32(b []byte, i *int32) (output []byte, n int, err error) {
	if len(b) < Encode32Size {
		return b, 0, ErrInputSize
	}
	if b[0] != 0x02 {
		return b, 0, ErrType
	}
	if i == nil {
		return b, 0, ErrNil
	}
	*i = int32(binary.BigEndian.Uint32(b[1:Encode32Size]))
	return b[Encode32Size:], Encode32Size, nil
}

const Encode64Size = 9

// EncodeInt64 encodes an int64 into out. It returns out, advanced by the data written, as well as the length of data written.
func EncodeInt64(i int64, out []byte) (output []byte, n int, err error) {
	if cap(out) < Encode64Size {
		return out, 0, ErrOutputSize
	}
	out = out[0:Encode64Size]
	out[0] = 0x03
	binary.BigEndian.PutUint64(out[1:Encode64Size], uint64(i))
	return out[Encode64Size:], Encode64Size, nil
}

// DecodeInt64 decodes an Int64 from in, it advances and returns in as well as the number of bytes read.
func DecodeInt64(in []byte, i *int64) (output []byte, n int, err error) {
	if len(in) < Encode64Size {
		return in, 0, ErrInputSize
	}
	if in[0] != 0x03 {
		return in, 0, ErrType
	}
	if i == nil {
		return in, 0, ErrNil
	}
	*i = int64(binary.BigEndian.Uint64(in[1:Encode64Size]))
	return in[Encode64Size:], Encode64Size, nil
}

// EncodeBytesSize returns the number of bytes d would occupy in encoded form.
func EncodeBytesSize(d []byte) int {
	if d == nil {
		return 5
	}
	return 5 + len(d)
}

// DecodeBytesSize returns the number of bytes the output buffer requires. It returns false if the entry cannot be decoded.
func DecodeBytesSize(in []byte) (int, bool) {
	if in[0] != 0x04 {
		return 0, false
	}
	if len(in) < 5 {
		return 0, false
	}
	l := int(binary.BigEndian.Uint32(in[1:5]))
	if l < 0 {
		return 0, false
	}
	if len(in) < 5+l {
		return 0, false
	}
	return l, true
}

// DecodeBytesSizeLimits returns the number of bytes the output buffer requires. It returns false if the entry cannot be decoded.
// It takes the given limits into account.
func DecodeBytesSizeLimits(in []byte, minSize, maxSize int) (int, bool) {
	size, ok := DecodeBytesSize(in)
	if !ok {
		return 0, false
	}
	if size < minSize {
		return 0, false
	}
	if maxSize > 0 && size > maxSize {
		return 0, false
	}
	return size, true
}

// EncodeBytes encodes a byte slice d into out. It returns out, advanced by the data written, as well as the length of data written.
// If out == nil, a zero length slice will be encoded.
func EncodeBytes(d, out []byte) (output []byte, n int, err error) {
	var l int
	size := EncodeBytesSize(d)
	if cap(out) < size {
		return out, 0, ErrOutputSize
	}
	out = out[0:size]
	out[0] = 0x04
	if d == nil {
		l = 0
	} else {
		l = len(d)
	}
	binary.BigEndian.PutUint32(out[1:5], uint32(l))
	if l > 0 {
		copy(out[5:size], d)
	}
	return out[size:], size, nil
}

// DecodeBytes decodes a byteslice from in into out. It returns out, advanced by the data written, as well as the length of data written.
// If out points to a nil slize, a new slize of the appropriate length will be created. If it points to a slice of zero length, data will
// be copied into that slize. It is an error if the input data is larger then the capacity. If out points to a slice of non-zero length,
// the input data must have exactly that size.
func DecodeBytes(in []byte, out *[]byte) (output []byte, n int, err error) {
	var x []byte
	if in[0] != 0x04 {
		return in, 0, ErrType
	}
	size, _ := DecodeBytesSize(in)
	if len(in) < size+5 {
		return in, 0, ErrInputSize
	}
	if *out == nil {
		x = make([]byte, size)
		*out = x
	}
	x = *out
	if size > 0 {
		if cap(x) < size {
			return in, 0, ErrSlizeExpectedLong
		}
		if len(x) > 0 && len(x) != size {
			return in, 0, ErrSlizeExpected
		}
		x = x[0:size]
		copy(x, in[5:5+size])
	}
	*out = x
	return in[5+size:], 5 + size, nil
}

// SlicePointer is a convenience function to Decode bytes into arrays.
func SlicePointer(d []byte) *[]byte {
	return &d
}

// EncodeSkip skips l bytes in the output.
func EncodeSkip(skip int, out []byte) (output []byte, n int, err error) {
	if cap(out) < skip {
		return out, 0, ErrOutputSize
	}
	out = out[0:skip]
	return out[skip:], skip, nil
}

// DecodeSkip decodes a skip in the output.
func DecodeSkip(in []byte, skip int) (output []byte, n int, err error) {
	if len(in) < skip {
		return in, 0, ErrInputSize
	}
	return in[skip:], skip, nil
}

// EncodeSize returns the necessary size of a byteslize to contain the given input values.
// Only *int16, int16, *int32, int32, *int64, int64, []byte and *[]byte are supported for encoding.
// int values are considered to be skip instructions.
func EncodeSize(v ...interface{}) (int, error) {
	var s int
	for _, vi := range v {
		switch e := vi.(type) {
		case *int16:
			s += Encode16Size
		case int16:
			s += Encode16Size
		case *int32:
			s += Encode32Size
		case int32:
			s += Encode32Size
		case *int64:
			s += Encode64Size
		case int64:
			s += Encode64Size
		case *[]byte:
			s += EncodeBytesSize(*e)
		case []byte:
			s += EncodeBytesSize(e)
		case int:
			s += e
		default:
			return 0, ErrType
		}
	}
	return s, nil
}

// Encode the interfaces to out. If out == nil, a new slice will be allocated.
// Only *int16, int16, *int32, int32, *int64, int64, []byte and *[]byte are supported for encoding.
// int values are considered to be skip instructions.
func Encode(out []byte, v ...interface{}) ([]byte, error) {
	var err error
	var s, sc, n int
	if out == nil {
		s, err = EncodeSize(v...)
		if err != nil {
			return nil, err
		}
		out = make([]byte, 0, s)
	}
	orig := out
	for _, vi := range v {
		switch e := vi.(type) {
		case *int16:
			if out, n, err = EncodeInt16(*e, out); err != nil {
				return nil, err
			}
		case int16:
			if out, n, err = EncodeInt16(e, out); err != nil {
				return nil, err
			}
		case *int32:
			if out, n, err = EncodeInt32(*e, out); err != nil {
				return nil, err
			}
		case int32:
			if out, n, err = EncodeInt32(e, out); err != nil {
				return nil, err
			}
		case *int64:
			if out, n, err = EncodeInt64(*e, out); err != nil {
				return nil, err
			}
		case int64:
			if out, n, err = EncodeInt64(e, out); err != nil {
				return nil, err
			}
		case *[]byte:
			if out, n, err = EncodeBytes(*e, out); err != nil {
				return nil, err
			}
		case []byte:
			if out, n, err = EncodeBytes(e, out); err != nil {
				return nil, err
			}
		case int:
			if out, n, err = EncodeSkip(e, out); err != nil {
				return nil, err
			}
		default:
			return nil, ErrType
		}
		sc += n
	}
	return orig[0:sc], nil
}

// Decode in into v. Returns the remainder.
// Only *int16, *int32, *int64, []byte and *[]byte are supported for encoding.
// int values are considered to be skip instructions.
func Decode(in []byte, v ...interface{}) ([]byte, error) {
	var err error
	for _, vi := range v {
		switch e := vi.(type) {
		case *int16:
			if in, _, err = DecodeInt16(in, e); err != nil {
				return in, err
			}
		case *int32:
			if in, _, err = DecodeInt32(in, e); err != nil {
				return in, err
			}
		case *int64:
			if in, _, err = DecodeInt64(in, e); err != nil {
				return in, err
			}
		case *[]byte:
			if in, _, err = DecodeBytes(in, e); err != nil {
				return in, err
			}
		case int:
			if in, _, err = DecodeSkip(in, e); err != nil {
				return in, err
			}
		default:
			return in, ErrType
		}
	}
	return in, nil
}

// DescribeStruct returns a slice of pointers that can be used to encode and decode
// the struct v which must be given as a pointer. For production code this function
// should not be used since it uses reflection.
func DescribeStruct(v interface{}) []interface{} {
	str := reflect.ValueOf(v).Elem()
	values := make([]interface{}, 0, str.NumField())
	for i := 0; i < str.NumField(); i++ {
		if !str.Field(i).CanInterface() {
			continue
		}
		switch str.Field(i).Interface().(type) {
		case int16:
			values = append(values, str.Field(i).Addr().Interface().(*int16))
		case int32:
			values = append(values, str.Field(i).Addr().Interface().(*int32))
		case int64:
			values = append(values, str.Field(i).Addr().Interface().(*int64))
		case []byte:
			values = append(values, str.Field(i).Addr().Interface().(*[]byte))
		}
	}
	return values
}
