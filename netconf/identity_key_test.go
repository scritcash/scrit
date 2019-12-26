package netconf

import (
	"encoding/hex"
	"testing"
)

const (
	publicKey             = "edf9cbecc47b3d4403f81b356529cf38d7dac58698df46899aed8cb40e5adcb4"
	marshalledIdentityKey = "ed25519-7fnL7MR7PUQD-Bs1ZSnPONfaxYaY30aJmu2MtA5a3LQ"
)

var identityKey IdentityKey

func init() {
	var err error
	identityKey.SigAlgo = "ed25519"
	identityKey.PubKey, err = hex.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}
}

func TestIdentityKeyMarshal(t *testing.T) {
	ik := identityKey.Marshal()
	if ik != marshalledIdentityKey {
		t.Errorf("identityKey.Marshal() == %s != %s", ik, marshalledIdentityKey)
	}
}
