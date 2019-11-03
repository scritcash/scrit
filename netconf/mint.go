package netconf

// Mint defines the key list of a single mint for all epochs.
type Mint struct {
	Epochs []MintEpoch // corresponding to global epochs
}

// MintEpoch defines the key list of a single mint for a single epoch and
// where to reach the mint.
type MintEpoch struct {
	URLs        []string     // how to each the mint
	SignStart   string       // start of signing epoch
	SignEnd     string       // end of signing epoch
	ValidateEnd string       // end of validation epoch
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
