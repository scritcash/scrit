package binencode

import "encoding/binary"

// SetType can be used to set the type indicator of a marshalled type. It is always encoded in the first two bytes.
func SetType(d []byte, dataType uint16) error {
	if len(d) < 2 {
		return ErrOutputSize
	}
	binary.BigEndian.PutUint16(d[0:2], dataType)
	return nil
}

// GetType can be used to get the type indicator of a marshalled type. It is always encoded in the first two bytes.
func GetType(d []byte) (dataType uint16, err error) {
	if len(d) < 2 {
		return 0, ErrInputSize
	}
	return binary.BigEndian.Uint16(d[0:2]), nil
}

// GetTypeExpect can be used to test the type indicator of a marshalled type. It is always encoded in the first two bytes.
func GetTypeExpect(d []byte, dataType uint16) (err error) {
	t, err := GetType(d)
	if len(d) < 2 {
		return ErrInputSize
	}
	if t != dataType {
		return ErrType
	}
	return nil
}
