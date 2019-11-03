package netconf

import (
	"time"
)

// NetworkEpoch globally defines a verification epoch (signing plus validation
// epoch) on the network.
type NetworkEpoch struct {
	M                       uint64           // the quorum
	N                       uint64           // total number of mints
	SignStart               time.Time        // start of signing epoch
	SignEnd                 time.Time        // end of signing epoch
	ValidateEnd             time.Time        // end of validation epoch
	MintsAdded              []IdentityKey    // mints added in this epoch
	MintsRemoved            []IdentityKey    // mints removed in this epoch
	MintsReplaced           []KeyReplacement // mints replaced in this epoch
	DBCTypesAdded           []DBCType        // DBC types added in this epoch
	DBCTypesRemoved         []DBCType        // DBC types removed in this epoch
	MonetarySupplyIncrease  []Note           // new notes to print
	MonetarySupplyReduction []Note           // TODO: define burn process
}

// IdentityKey defines a mint identity key.
type IdentityKey struct {
	SigAlgo string // signature algorithm
	PubKey  []byte // public key
}

// KeyReplacement defines a mint identity key replacement.
type KeyReplacement struct {
	NewKey    IdentityKey // the new identity key
	OldKey    IdentityKey // the replaced identity key
	Signature string      // of new key by replaced key
}

// DBCType defines a DBC type.
type DBCType struct {
	Currency string // the DBC currency
	Amount   uint64 // per DBC
}

// Note defines newly printed or burned DBCs.
type Note struct {
	Random         [16]byte // nonce
	Quantity       uint64   // number of DBCs
	Currency       string   // DBC currency
	Amount         uint64   // per DBC
	ReceiverPubKey []byte   // recipient
}

// Validate the network epoch.
func (epoch *NetworkEpoch) Validate() error {
	// m > 0
	if epoch.M == 0 {
		return ErrZeroM
	}
	// n > 0
	if epoch.N == 0 {
		return ErrZeroN
	}
	// m <= n
	if epoch.M > epoch.N {
		return ErrMGreaterN
	}
	// m > n/2
	if epoch.M <= epoch.N/2 {
		return ErrQuorumTooSmall
	}

	// sign epoch start < sign epoch end
	/*
		if !epoch.SignStart.Before(epoch.SignEnd) {
			return ErrSignEpochStartAfterEnd
		}
	*/
	// sign epoch end < validation epoch end
	/*
		if !epoch.Sign.End(epoch.ValidateEnd) {
			return ErrSignEpochAfterValidate
		}
	*/

	return nil
}
