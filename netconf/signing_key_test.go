package netconf

import (
	"testing"
)

func TestSigningKey(t *testing.T) {
	_, err := NewSigningKey("EUR", 100000000)
	if err != nil {
		t.Error(err)
	}
}
