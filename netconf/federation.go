package netconf

import (
	"path/filepath"
)

// A Federation of Scrit mints.
type Federation struct {
	network *Network
	mints   map[string]*Mint
}

// LoadFederation loads a Scrit mint federation configuration from the given
// directory and validates it.
func LoadFederation(dir string) (*Federation, error) {
	var f Federation
	filename := filepath.Join(dir, DefNetConfFile)
	n, err := LoadNetwork(filename)
	if err != nil {
		return nil, err
	}
	if err := n.Validate(); err != nil {
		return nil, err
	}
	f.network = n

	return &f, nil
}
