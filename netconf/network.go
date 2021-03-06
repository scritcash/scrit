package netconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/frankbraun/codechain/util/file"
)

// Network defines a Scrit network.
type Network struct {
	NetworkEpochs []NetworkEpoch // global list of signing epochs
}

// NewNetwork creates a new network configuration and returns the Network
// struct.
func NewNetwork(
	m, n uint64,
	signStart, signEnd, validateEnd time.Time,
	mintIdentities []IdentityKey,
) *Network {
	var network Network
	network.NetworkEpochs = []NetworkEpoch{
		{
			QuorumM:        m,
			NumberOfMintsN: n,
			SignStart:      signStart,
			SignEnd:        signEnd,
			ValidateEnd:    validateEnd,
			MintsAdded:     mintIdentities,
		},
	}
	return &network
}

// LoadNetwork loads a network configuration from filename and return
// the Network struct.
func LoadNetwork(filename string) (*Network, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var n Network
	if err := json.Unmarshal(data, &n); err != nil {
		return nil, err
	}
	return &n, err
}

// Validate the network configuration.
func (n *Network) Validate() error {
	// validate each network epoch
	for _, e := range n.NetworkEpochs {
		if err := e.Validate(); err != nil {
			return err
		}
	}
	// validate network epoch transitions
	for i := 1; i < len(n.NetworkEpochs); i++ {
		// sign end i-1 == sign start i
		if n.NetworkEpochs[i-1].SignEnd != n.NetworkEpochs[i].SignStart {
			return ErrSignEpochWrongBoundaries
		}
		// validation end i-1 <= sign end i
		if n.NetworkEpochs[i-1].ValidateEnd.After(n.NetworkEpochs[i].SignEnd) {
			return ErrValidationLongerThanNextSigning
		}
	}

	// validate Mints
	if err := n.MintsValidate(); err != nil {
		return err
	}

	// validate DBC types
	if err := n.DBCTypesValidate(); err != nil {
		return err
	}

	return nil
}

// Marshal network as string.
func (n *Network) Marshal() string {
	jsn, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		panic(err) // should never happen
	}
	return string(jsn)
}

// Save network to filename. If filename exists already it will be
// overwritten!
func (n *Network) Save(filename string) error {
	jsn, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return err
	}
	exists, err := file.Exists(filename)
	if err != nil {
		return err
	}
	if exists {
		if err := os.Rename(filename, filename+".bac"); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(filename, jsn, 0755); err != nil {
		return err
	}
	if exists {
		return os.Remove(filename + ".bac")
	}
	return nil
}

// HasFuture ensures that the network has an epoch which starts in the future.
func (n *Network) HasFuture() error {
	e := n.NetworkEpochs[len(n.NetworkEpochs)-1]
	if !e.SignStart.After(time.Now().UTC()) {
		return ErrNoFuture
	}
	return nil
}

// Mints returns a map of all mints in the network in the future.
func (n *Network) Mints() map[string]bool {
	mints := make(map[string]bool)
	for _, e := range n.NetworkEpochs {
		for _, add := range e.MintsAdded {
			mints[add.MarshalID()] = true
		}
		for _, remove := range e.MintsRemoved {
			delete(mints, remove.MarshalID())
		}
		for _, replace := range e.MintsReplaced {
			delete(mints, replace.OldKey.MarshalID())
			mints[replace.NewKey.MarshalID()] = true
		}
	}
	return mints
}

// AllMints returns a map of all mints that were ever or will ever be part of
// the network.
func (n *Network) AllMints() map[string]bool {
	mints := make(map[string]bool)
	for _, e := range n.NetworkEpochs {
		for _, add := range e.MintsAdded {
			mints[add.MarshalID()] = true
		}
		for _, replace := range e.MintsReplaced {
			mints[replace.NewKey.MarshalID()] = true
		}
	}
	return mints
}

// CurrentMints returns a map of all mints in the network at the current time
func (n *Network) CurrentMints() (map[string]bool, error) {
	c, err := n.CurrentEpoch()
	if err != nil {
		return nil, err
	}
	mints := make(map[string]bool)
	for i := 0; i <= c; i++ {
		e := n.NetworkEpochs[i]
		for _, add := range e.MintsAdded {
			mints[add.MarshalID()] = true
		}
		for _, remove := range e.MintsRemoved {
			delete(mints, remove.MarshalID())
		}
		for _, replace := range e.MintsReplaced {
			delete(mints, replace.OldKey.MarshalID())
			mints[replace.NewKey.MarshalID()] = true
		}
	}
	return mints, nil
}

// MintsValidate validates the mint types.
func (n *Network) MintsValidate() error {
	mints := make(map[string]bool)
	for _, e := range n.NetworkEpochs {
		// make sure the MintsAdded, MintsRemoved, and MintsReplaced sets are
		// disjunct
		if err := e.MintsDisjunct(); err != nil {
			return err
		}
		for _, add := range e.MintsAdded {
			// make sure we do not add an exisiting mint
			id := add.MarshalID()
			if mints[id] {
				return fmt.Errorf("netconf: mint already added: %v", id)
			}
			mints[id] = true
		}
		for _, remove := range e.MintsRemoved {
			id := remove.MarshalID()
			// make sure the mint to delete is actually there
			_, present := mints[id]
			if !present {
				return fmt.Errorf("netconf: mint to remove not added: %v", id)
			}
			delete(mints, id)
		}
		for _, replace := range e.MintsReplaced {
			if err := replace.Verify(); err != nil {
				return err
			}
			oldID := replace.OldKey.MarshalID()
			newID := replace.NewKey.MarshalID()
			// make sure the mint to replace is actually there
			_, present := mints[oldID]
			if !present {
				return fmt.Errorf("netconf: mint to replace not added: %v", oldID)
			}
			delete(mints, oldID)
			// make sure we do not replace to an exisiting mint
			if mints[newID] {
				return fmt.Errorf("netconf: mint to replace to already added: %v", newID)
			}
			mints[newID] = true
		}
	}
	return nil
}

// MintAdd adds the mint identity key to the network.
// Low-level function without error checking!
func (n *Network) MintAdd(key *IdentityKey) {
	n.NetworkEpochs[len(n.NetworkEpochs)-1].MintsAdded =
		append(n.NetworkEpochs[len(n.NetworkEpochs)-1].MintsAdded, *key)
}

// MintRemove removes the mint identity key from the network.
// Low-level function without error checking!
func (n *Network) MintRemove(key *IdentityKey) {
	n.NetworkEpochs[len(n.NetworkEpochs)-1].MintsRemoved =
		append(n.NetworkEpochs[len(n.NetworkEpochs)-1].MintsRemoved, *key)
}

// MintReplace replaces the old mint identity key with the new mint identity
// in the network.
// Low-level function without error checking!
func (n *Network) MintReplace(r *KeyReplacement) {
	n.NetworkEpochs[len(n.NetworkEpochs)-1].MintsReplaced =
		append(n.NetworkEpochs[len(n.NetworkEpochs)-1].MintsReplaced, *r)
}

// DBCTypes returns a map of all DBCTypes in the network.
func (n *Network) DBCTypes() map[DBCType]bool {
	dbcTypes := make(map[DBCType]bool)
	for _, e := range n.NetworkEpochs {
		for _, add := range e.DBCTypesAdded {
			dbcTypes[add] = true
		}
		for _, remove := range e.DBCTypesRemoved {
			delete(dbcTypes, remove)
		}
	}
	return dbcTypes
}

// DBCTypesValidate validates the DBC types.
func (n *Network) DBCTypesValidate() error {
	dbcTypes := make(map[DBCType]bool)
	for _, e := range n.NetworkEpochs {
		// make sure the DBCTypesAdded and DBCTypesRemoved sets are disjunct
		if err := e.DBCTypesDisjunct(); err != nil {
			return err
		}
		for _, add := range e.DBCTypesAdded {
			// make sure we do not add an exisiting DBC type
			if dbcTypes[add] {
				return fmt.Errorf("netconf: DBC type already defined: %v", add)
			}
			dbcTypes[add] = true
		}
		for _, remove := range e.DBCTypesRemoved {
			// make sure the type to delete is actually there
			_, present := dbcTypes[remove]
			if !present {
				return fmt.Errorf("netconf: DBC type not defined: %v", remove)
			}
			delete(dbcTypes, remove)
		}
	}
	return nil
}

// DBCTypeAdd adds the DBC type to the network.
// Low-level function without error checking!
func (n *Network) DBCTypeAdd(dt DBCType) {
	n.NetworkEpochs[len(n.NetworkEpochs)-1].DBCTypesAdded =
		append(n.NetworkEpochs[len(n.NetworkEpochs)-1].DBCTypesAdded, dt)
}

// DBCTypeRemove removes the DBC type from the network.
// Low-level function without error checking!
func (n *Network) DBCTypeRemove(dt DBCType) {
	n.NetworkEpochs[len(n.NetworkEpochs)-1].DBCTypesRemoved =
		append(n.NetworkEpochs[len(n.NetworkEpochs)-1].DBCTypesRemoved, dt)
}

// EpochAdd adds another epoch with the given signing and validation period to
// the network.
// Low-level function without error checking!
func (n *Network) EpochAdd(signingPeriod, validationPeriod time.Duration) {
	lastEpoch := n.NetworkEpochs[len(n.NetworkEpochs)-1]
	var newEpoch NetworkEpoch
	newEpoch.QuorumM = lastEpoch.QuorumM
	newEpoch.NumberOfMintsN = lastEpoch.NumberOfMintsN
	newEpoch.SignStart = lastEpoch.SignEnd
	newEpoch.SignEnd = newEpoch.SignStart.Add(signingPeriod)
	newEpoch.ValidateEnd = newEpoch.SignEnd.Add(validationPeriod)
	n.NetworkEpochs = append(n.NetworkEpochs, newEpoch)
}

// SetQuorum sets the quorum for the last epoch.
// Low-level function without error checking!
func (n *Network) SetQuorum(m uint64) {
	n.NetworkEpochs[len(n.NetworkEpochs)-1].QuorumM = m
}

// CurrentEpoch returns the current signing epoch number or an error if no
// such epoch exists.
func (n *Network) CurrentEpoch() (int, error) {
	// make sure we still have a valid signing epoch
	i := len(n.NetworkEpochs) - 1
	now := time.Now().UTC()
	if now.After(n.NetworkEpochs[i].SignEnd) {
		return 0, errors.New("netconf: no valid signing epoch found")
	}
	// determine current epoch
	for ; i >= 0; i-- {
		if now.After(n.NetworkEpochs[i].SignStart) {
			break
		}
	}
	if i < 0 {
		i = 0
	}
	return i, nil
}
