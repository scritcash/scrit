package mintcom

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	_ "crypto/sha256" // USed hashing algorithm.
	"errors"
	"time"
)

// Errors
var (
	ErrShortRead = errors.New("mitcom: Short read from reader")
)

// Package markers
const (
	// PkgTypeCommitment is a commitment message.
	PkgTypeCommitment = byte(0x01)
)

// PublicKeyLookup is a function type that returns the corresponding public key for a mintID, or nil if the key cannot be found.
type PublicKeyLookup func(mintID uint64) *[PublicKeySize]byte

// Now is the source of current unixtime for this package.
var Now = func() uint64 { return uint64(time.Now().Unix()) }

// Rand is the source of random bytes for this package.
var Rand = rand.Reader

// HashAlgo is the hashing algorithm used. Change requires updating HashSize and imports.
const HashAlgo = crypto.SHA256

// HashSize is the size of hashes in bytes.
const HashSize = 32

// Hash returns the hash of i.
func Hash(i []byte) [HashSize]byte {
	r := new([HashSize]byte)
	h := HashAlgo.New()
	h.Write(i)
	h.Sum(r[0:0])
	return *r
}

// HMAC returns an HMAC for msg using key.
func HMAC(msg, key []byte) [HashSize]byte {
	r := new([HashSize]byte)
	h := hmac.New(HashAlgo.New, key)
	h.Write(msg)
	h.Sum(r[0:0])
	return *r
}

// RandomBytes fills d with random bytes.
func RandomBytes(d []byte) error {
	n, err := Rand.Read(d)
	if err != nil {
		return err
	}
	if n < len(d) {
		return ErrShortRead
	}
	return nil
}
