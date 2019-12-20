package netconf

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

// Mint defines the key list of a single mint for all epochs.
type Mint struct {
	Description     string      // description of mint (name)
	MintIdentityKey IdentityKey // identity key of mint
	MintEpochs      []MintEpoch // corresponding to global epochs
}

// MintEpoch defines the key list of a single mint for a single epoch and
// where to reach the mint.
type MintEpoch struct {
	URLs        []string     // how to reach the mint
	SignStart   time.Time    // start of signing epoch
	SignEnd     time.Time    // end of signing epoch
	ValidateEnd time.Time    // end of validation epoch
	KeyList     []SigningKey // the key list
}

// SigningKey defines an entry in the key list.
type SigningKey struct {
	Currency          string // the currency this key signs
	Amount            uint64 // the amount this key signs
	SigAlgo           string // signature algorithm
	PubKey            []byte // public key
	SelfSignature     []byte // self signature
	IdentitySignature []byte // signature by identity key
}

// Load a mint configuration from filename and return the Mint struct.
func LoadMint(filename string) (*Mint, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var mint Mint
	if err := json.Unmarshal(data, &mint); err != nil {
		return nil, err
	}
	return &mint, err
}

// Validate the mint configuration.
func (mint *Mint) Validate() error {
	// TODO
	return nil
}
