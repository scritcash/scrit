package netconf

import (
	"crypto/ed25519"
	"crypto/rand"
)

// SigningKey defines an entry in the key list.
type SigningKey struct {
	Currency string             // the currency this key signs, usually ISO 4217 codes
	Amount   uint64             // the amount this key signs, 8 digits after the dot
	SigAlgo  string             // signature algorithm
	PubKey   []byte             // public key
	privKey  ed25519.PrivateKey // private key
}

// NewSigningKey generates a new signing key.
func NewSigningKey(
	currency string,
	amount uint64,
	ik *IdentityKey,
) (*SigningKey, error) {
	var sk SigningKey
	sk.Currency = currency
	sk.Amount = amount
	sk.SigAlgo = "ed25519" // TODO
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	sk.PubKey = pubKey
	sk.privKey = privKey
	return &sk, nil
}
