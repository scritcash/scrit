package netconf

import (
	"encoding/base64"
)

// IdentityKey defines a mint identity key.
type IdentityKey struct {
	SigAlgo string // signature algorithm
	PubKey  []byte // public key
}

func (ik *IdentityKey) Marshal() string {
	return ik.SigAlgo + "-" + base64.RawURLEncoding.EncodeToString(ik.PubKey)
}
