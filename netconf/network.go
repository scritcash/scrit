package netconf

import (
	"encoding/json"
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
		if err := os.Rename(filename, "."+filename); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(filename, jsn, 0755); err != nil {
		return err
	}
	if exists {
		return os.Remove("." + filename)
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

// DBCTypes returns a map of all DBCTypes in network.
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
