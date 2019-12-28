package netconf

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

// Marshal identity key.
func (ik *IdentityKey) Marshal() string {
	return ik.SigAlgo + "-" + base64.RawURLEncoding.EncodeToString(ik.PubKey)
}

// MarshalJSON ik as JSON string.
func (ik *IdentityKey) MarshalJSON() string {
	jsn, err := json.MarshalIndent(ik, "", "  ")
	if err != nil {
		panic(err) // should never happen
	}
	return string(jsn)
}
