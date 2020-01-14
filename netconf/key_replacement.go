package netconf

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

// KeyReplacement defines a mint identity key replacement.
type KeyReplacement struct {
	NewKey    IdentityKey // the new identity key
	OldKey    IdentityKey // the replaced identity key
	Signature string      // of new key by replaced key
}

// NewKeyReplacement returns a new key replacement from oldKey to newKey.
// The signature is from oldKey over newKey.
func NewKeyReplacement(newKey, oldKey *IdentityKey, sig string) *KeyReplacement {
	return &KeyReplacement{
		NewKey:    *newKey,
		OldKey:    *oldKey,
		Signature: sig,
	}
}

// Verify the signature of the given key replacement
func (r *KeyReplacement) Verify() error {
	sig, err := base64.RawURLEncoding.DecodeString(r.Signature)
	if err != nil {
		return err
	}
	if !ed25519.Verify(r.OldKey.PubKey, []byte(r.NewKey.MarshalID()), sig) {
		return fmt.Errorf("netconf: signature '%s' does not verify", sig)
	}
	return nil
}
