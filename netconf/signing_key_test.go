package netconf

import (
	"testing"
)

func TestSigningKey(t *testing.T) {
	ik, err := NewIdentityKey()
	if err != nil {
		t.Error(err)
	}
	sk, err := NewSigningKey("EUR", 100000000, ik)
	if err != nil {
		t.Error(err)
	}
	if err := sk.Verify(ik); err != nil {
		t.Error(err)
	}
}
