package mintcom

import (
	"bytes"
	"crypto"
	"io"
	"testing"
)

// TestIntegrity checks the integrity of the constants in this package.
func TestIntegrity(t *testing.T) {
	// Verify that the hash algorithm used has the correct HashSize.
	h := HashAlgo.New()
	if h.Size() != HashSize {
		t.Fatalf("HashSize (%d) and HashAlgo (%d) do not match.", HashSize, h.Size())
	}
}

// TestHash verifies Hash()
func TestHash(t *testing.T) {
	var td []byte
	ts := []byte("test1")
	switch HashAlgo {
	case crypto.SHA256:
		td = []byte{0x1b, 0x4f, 0x0e, 0x98, 0x51, 0x97, 0x19, 0x98, 0xe7, 0x32, 0x07, 0x85, 0x44, 0xc9, 0x6b, 0x36, 0xc3, 0xd0, 0x1c, 0xed, 0xf7, 0xca, 0xa3, 0x32, 0x35, 0x9d, 0x6f, 0x1d, 0x83, 0x56, 0x70, 0x14}
	default:
		t.Fatal("No test for HashAlgo.")
	}
	o := Hash(ts)
	if !bytes.Equal(o[:], td) {
		t.Error("Hash() returns wrong value.")
	}
}

// TestHMAC verifies HMAC()
func TestHMAC(t *testing.T) {
	var td []byte
	tk := []byte("TestKey")
	tm := []byte("TestMsg")
	switch HashAlgo {
	case crypto.SHA256:
		td = []byte{0x34, 0x56, 0xd1, 0x3f, 0x33, 0x5a, 0xf5, 0x6c, 0x71, 0xc5, 0x45, 0x87, 0x35, 0x90, 0xb3, 0xbc, 0x6c, 0x4e, 0x0a, 0x81, 0x20, 0x34, 0x8a, 0xcd, 0xdd, 0xba, 0xf4, 0x3b, 0xd0, 0x55, 0x5c, 0xd6}
	default:
		t.Fatal("No test for HashAlgo.")
	}
	o := HMAC(tm, tk)
	if !bytes.Equal(o[:], td) {
		t.Error("HMAC() returns wrong value.")
	}
}

//===================

type ErrorReader struct{}

func (e ErrorReader) Read(d []byte) (int, error) {
	return len(d), io.ErrUnexpectedEOF
}

type ShortReader struct{}

func (e ShortReader) Read(d []byte) (int, error) {
	return len(d) - 1, nil
}

// TestRandomBytes verifies RandomBytes
func TestRandomBytes(t *testing.T) {
	td1 := []byte("0123456789")
	td2 := []byte("0123456789")
	// Test normal output.
	if err := RandomBytes(td2); err != nil {
		t.Errorf("RandomBytes returned unexpected error: %s", err)
	}
	if bytes.Equal(td1, td2) {
		t.Error("RandomBytes did not change output slice.")
	}
	// Prepare readers with deliberate errors.
	origRand := Rand
	defer func() { Rand = origRand }()

	// Reader returns error.
	Rand = ErrorReader{}
	if err := RandomBytes(td2); err != io.ErrUnexpectedEOF {
		t.Error("RandomBytes did not return reader error")
	}

	// Reader reads short.
	Rand = ShortReader{}
	if err := RandomBytes(td2); err != ErrShortRead {
		t.Error("RandomBytes did not recognize short read")
	}

}
