package netconf

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
)

// IdentityKey defines a mint identity key.
type IdentityKey struct {
	SigAlgo string             // signature algorithm
	PubKey  []byte             // public key
	privKey ed25519.PrivateKey // private key
}

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

func (ik *IdentityKey) Marshal() string {
	return ik.SigAlgo + "-" + base64.RawURLEncoding.EncodeToString(ik.PubKey)
}
