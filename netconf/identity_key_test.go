package netconf

import (
	"encoding/hex"
	"testing"
)

const (
	publicKey             = "38ff60e128224b4eb708a49de1908d7d61b447a8228459b216ac49c209910295"
	privateKey            = "4b75408909900de74555aa57bb667ab36e18d242f4962b317a0fedc2083891c038ff60e128224b4eb708a49de1908d7d61b447a8228459b216ac49c209910295"
	marshalledIdentityKey = "ed25519-OP9g4SgiS063CKSd4ZCNfWG0R6gihFmyFqxJwgmRApU"
)

var identityKey IdentityKey

func init() {
	var err error
	identityKey.SigAlgo = "ed25519"
	identityKey.PubKey, err = hex.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}
	identityKey.privKey, err = hex.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}
}

func TestIdentityKeyMarshal(t *testing.T) {
	iks := identityKey.MarshalID()
	if iks != marshalledIdentityKey {
		t.Errorf("identityKey.MarshalID() == %s != %s", iks, marshalledIdentityKey)
	}
	ik, err := ParseIdentityKey(iks)
	if err != nil {
		t.Error(err)
	}
	iks2 := ik.MarshalID()
	if iks2 != iks {
		t.Errorf("ik.MarshalID() == %s != %s", iks2, iks)
	}
}
