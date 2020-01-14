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
	Amount   uint64 // the amount per DBC, last 8 digits are decimal places
}

// Validate the network epoch.
func (e *NetworkEpoch) Validate() error {
	// m > 0
	if e.QuorumM == 0 {
		return ErrZeroM
	}
	// n > 0
	if e.NumberOfMintsN == 0 {
		return ErrZeroN
	}
	// m <= n
	if e.QuorumM > e.NumberOfMintsN {
		return ErrMGreaterN
	}
	// m > n/2
	if e.QuorumM <= e.NumberOfMintsN/2 {
		return ErrQuorumTooSmall
	}

	// sign epoch start < sign epoch end
	if !e.SignStart.Before(e.SignEnd) {
		return ErrSignEpochStartNotBeforeSignEnd
	}
	// sign epoch end < validation epoch end
	if !e.SignEnd.Before(e.ValidateEnd) {
		return ErrSignEpochEndNotBeforeValidateEnd
	}

	return nil
}

// MintsDisjunct make sure the MintsAdded, MintsRemoved, and MintsReplaced
// sets are disjunct.
func (e *NetworkEpoch) MintsDisjunct() error {
	addedMints := make(map[string]bool)
	removedMints := make(map[string]bool)
	replacedMints := make(map[string]bool)
	// fill maps
	for _, add := range e.MintsAdded {
		addedMints[add.MarshalID()] = true
	}
	for _, remove := range e.MintsRemoved {
		removedMints[remove.MarshalID()] = true
	}
	for _, replace := range e.MintsReplaced {
		replacedMints[replace.OldKey.MarshalID()] = true
	}
	// check
	for _, replace := range e.MintsReplaced {
		newID := replace.NewKey.MarshalID()
		oldID := replace.OldKey.MarshalID()
		if addedMints[newID] || removedMints[newID] || replacedMints[newID] {
			return ErrMintsOverlap
		}
		if addedMints[oldID] || removedMints[oldID] {
			return ErrMintsOverlap
		}
	}
	for _, remove := range e.MintsRemoved {
		if addedMints[remove.MarshalID()] {
			return ErrMintsOverlap
		}
	}
	return nil
}

// DBCTypesDisjunct makes sure the DBCTypesAdded and DBCTypesRemoved sets from the epoch are disjunct.
func (e *NetworkEpoch) DBCTypesDisjunct() error {
	dbcTypes := make(map[DBCType]bool)
	for _, add := range e.DBCTypesAdded {
		dbcTypes[add] = true
	}
	for _, remove := range e.DBCTypesRemoved {
		if dbcTypes[remove] {
			return ErrDBCTypesOverlap
		}
	}
	return nil
}
