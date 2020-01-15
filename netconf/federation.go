package netconf

import (
	"fmt"
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
	fmt.Printf("loading '%s'\n", filename)
	n, err := LoadNetwork(filename)
	if err != nil {
		return nil, err
	}
	fmt.Printf("validate '%s'\n", filename)
	if err := n.Validate(); err != nil {
		return nil, err
	}
	f.network = n
	f.mints = make(map[string]*Mint)

	// TODO: the following is not correct in all situations, we have to validate
	// in the current signing period (and maybe later ones)

	mints := n.Mints()
	for mn := range mints {
		filename := filepath.Join(dir, DefMintDir, mn+".json")
		fmt.Printf("loading '%s'\n", filename)
		m, err := LoadMint(filename)
		if err != nil {
			return nil, err
		}
		fmt.Printf("validate '%s'\n", filename)
		if err := m.Validate(); err != nil {
			return nil, err
		}
		f.mints[mn] = m
	}
	return &f, nil
}
