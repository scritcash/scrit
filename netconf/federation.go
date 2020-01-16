package netconf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// A Federation of Scrit mints.
type Federation struct {
	Network *Network         // the federation network
	Mints   map[string]*Mint // all mints in the network
}

func (f *Federation) validate() error {
	// for every mint we make sure the mint epochs match with the network epoch
	for _, m := range f.Mints {
		for i, e := range m.MintEpochs {
			if e.SignStart != f.Network.NetworkEpochs[i].SignStart {
				return errors.New("netconf: signing start mismatch")
			}
			if e.SignEnd != f.Network.NetworkEpochs[i].SignEnd {
				return errors.New("netconf: signing end mismatch")
			}
			if e.ValidateEnd != f.Network.NetworkEpochs[i].ValidateEnd {
				return errors.New("netconf: validation end mismatch")
			}
		}
	}

	// now we make sure that in the present we have enough mint epochs (quorum)
	i, err := f.Network.CurrentEpoch()
	if err != nil {
		return err
	}
	var q uint64
	for _, m := range f.Mints {
		if len(m.MintEpochs) > i {
			q++
		}
	}
	if q < f.Network.NetworkEpochs[i].QuorumM {
		return errors.New("netconf: not enough mints to reach quorum in present")
	}

	return nil
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
	f.Network = n
	f.Mints = make(map[string]*Mint)

	// we try to load all mints of the current signing epoch, but ignore errors
	// f.validate() later checks that we have enough mints available
	mints, err := n.CurrentMints()
	if err != nil {
		return nil, err
	}
	for mn := range mints {
		filename := filepath.Join(dir, DefMintDir, mn+".json")
		m, err := LoadMint(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING loading '%s' failed: %s\n", filename, err)
			continue
		}
		if err := m.Validate(n); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING validating '%s' failed: %s\n", filename, err)
			continue
		}
		f.Mints[mn] = m
	}
	if err := f.validate(); err != nil {
		return nil, err
	}
	return &f, nil
}
