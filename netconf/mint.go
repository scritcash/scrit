package netconf

import (
	//"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/scritcash/scrit/binencode"
)

// Mint defines the key list of a single mint for all epochs and where to
// reach the mint.
type Mint struct {
	Description     string      // description of mint (name)
	MintIdentityKey IdentityKey // identity key of mint
	MintEpochs      []MintEpoch // corresponding to global epochs
	URLs            []string    // how to reach the mint
}

// MintEpoch defines the key list of a single mint for a single epoch.
type MintEpoch struct {
	SignStart         time.Time    // start of signing epoch
	SignEnd           time.Time    // end of signing epoch
	ValidateEnd       time.Time    // end of validation epoch
	KeyList           []SigningKey // the key list
	KeyListSignatures [][]byte     // signatures of key list (identity signature last)
}

// LoadMint loads  a mint configuration from filename and return the
// Mint struct.
func LoadMint(filename string) (*Mint, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var mint Mint
	if err := json.Unmarshal(data, &mint); err != nil {
		return nil, err
	}

	// TODO
	for _, me := range mint.MintEpochs {
		if err := me.Sign(&mint.MintIdentityKey); err != nil {
			return nil, err
		}
	}

	return &mint, err
}

// Sign mint epoch.
func (me *MintEpoch) Sign(ik *IdentityKey) error {
	encodingScheme := []interface{}{
		[]byte(ik.SigAlgo),
		ik.PubKey,
		me.SignStart.UTC().Unix(),
		me.SignEnd.UTC().Unix(),
		me.ValidateEnd.UTC().Unix(),
	}
	for _, k := range me.KeyList {
		encodingScheme = append(encodingScheme,
			[]byte(k.Currency),
			int64(k.Amount),
			[]byte(k.SigAlgo),
			k.PubKey,
		)
	}
	size, err := binencode.EncodeSize(encodingScheme...)
	if err != nil {
		return err
	}
	fmt.Printf("size=%d\n", size)
	buf := make([]byte, size)
	_, err = binencode.Encode(buf, encodingScheme...)
	if err != nil {
		return err
	}
	/*
		for _, k := range me.KeyList {
			sig := ed25519.Sign(k.privKey, enc)
			me.KeyListSignatures = append(me.KeyListSignatures, sig)
		}
		sig := ed25519.Sign(ik.privKey, enc)
		me.KeyListSignatures = append(me.KeyListSignatures, sig)
	*/
	return nil
}

// Verify mint epoch.
func (me *MintEpoch) Verify(ik *IdentityKey) error {
	// TODO
	return nil
}

// Validate the mint configuration.
func (mint *Mint) Validate() error {
	// TODO
	return nil
}
