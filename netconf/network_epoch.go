package netconf

import (
	"time"
)

// NetworkEpoch globally defines a verification epoch (signing plus validation
// epoch) on the network.
type NetworkEpoch struct {
	QuorumM         uint64           // the quorum
	NumberOfMintsN  uint64           // total number of mints
	SignStart       time.Time        // start of signing epoch
	SignEnd         time.Time        // end of signing epoch
	ValidateEnd     time.Time        // end of validation epoch
	MintsAdded      []IdentityKey    // mints added in this epoch
	MintsRemoved    []IdentityKey    // mints removed in this epoch
	MintsReplaced   []KeyReplacement // mints replaced in this epoch
	DBCTypesAdded   []DBCType        // DBC types added in this epoch
	DBCTypesRemoved []DBCType        // DBC types removed in this epoch
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

// Validate the network epoch.
func (epoch *NetworkEpoch) Validate() error {
	// m > 0
	if epoch.QuorumM == 0 {
		return ErrZeroM
	}
	// n > 0
	if epoch.NumberOfMintsN == 0 {
		return ErrZeroN
	}
	// m <= n
	if epoch.QuorumM > epoch.NumberOfMintsN {
		return ErrMGreaterN
	}
	// m > n/2
	if epoch.QuorumM <= epoch.NumberOfMintsN/2 {
		return ErrQuorumTooSmall
	}

	// sign epoch start < sign epoch end
	if !epoch.SignStart.Before(epoch.SignEnd) {
		return ErrSignEpochStartNotBeforeSignEnd
	}
	// sign epoch end < validation epoch end
	if !epoch.SignEnd.Before(epoch.ValidateEnd) {
		return ErrSignEpochEndNotBeforeValidateEnd
	}

	return nil
}
