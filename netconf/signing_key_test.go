package netconf

import (
	"testing"
)

func TestSigningKey(t *testing.T) {
	ik, err := NewIdentityKey()
	if err != nil {
		t.Error(err)
	}
	_, err = NewSigningKey("EUR", 100000000, ik)
	if err != nil {
		t.Error(err)
	}
}
