package netconf

import (
	"encoding/json"
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
