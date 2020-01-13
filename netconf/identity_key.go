package netconf

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// IdentityKey defines a mint identity key.
type IdentityKey struct {
	SigAlgo string             // signature algorithm
	PubKey  []byte             // public key
	privKey ed25519.PrivateKey // private key
}

// NewIdentityKey generates a new identity key.
func NewIdentityKey() (*IdentityKey, error) {
	var ik IdentityKey
	ik.SigAlgo = "ed25519" // TODO
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	ik.PubKey = pubKey
	ik.privKey = privKey
	return &ik, err
}

// ParseIdentityKey parses a mint identity key (one liner).
func ParseIdentityKey(iks string) (*IdentityKey, error) {
	var ik IdentityKey
	parts := strings.SplitN(iks, "-", 2)
	ik.SigAlgo = parts[0]
	pk, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("netconf: cannot parse identity key '%s': %s",
			iks, err)
	}
	ik.PubKey = pk
	return &ik, nil
}

// NewIdentityKeyEd25519Priv create a new Ed25519 identity key from the given
// Ed25519 private key.
func NewIdentityKeyEd25519Priv(privKey *[64]byte) *IdentityKey {
	var ik IdentityKey
	ik.SigAlgo = "ed25519"
	ik.PubKey = make([]byte, 32)
	copy(ik.PubKey, privKey[32:])
	ik.privKey = make([]byte, 64)
	copy(ik.privKey, privKey[:])
	return &ik
}

// MarshalID marshals identity key as ID.
func (ik *IdentityKey) MarshalID() string {
	return ik.SigAlgo + "-" + base64.RawURLEncoding.EncodeToString(ik.PubKey)
}

// Marshal ik as JSON string.
func (ik *IdentityKey) Marshal() string {
	jsn, err := json.MarshalIndent(ik, "", "  ")
	if err != nil {
		panic(err) // should never happen
	}
	return string(jsn)
}
