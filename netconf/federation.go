package netconf

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"
)

// A Federation of Scrit mints.
type Federation struct {
	network *Network
	mints   map[string]*Mint
}

func (f *Federation) validate() error {
	// for every mint we make sure the mint epochs match with the network epoch
	for _, m := range f.mints {
		for i, e := range m.MintEpochs {
			if e.SignStart != f.network.NetworkEpochs[i].SignStart {
				return errors.New("netconf: signing start mismatch")
			}
			if e.SignEnd != f.network.NetworkEpochs[i].SignEnd {
				return errors.New("netconf: signing end mismatch")
			}
			if e.ValidateEnd != f.network.NetworkEpochs[i].ValidateEnd {
				return errors.New("netconf: validation end mismatch")
			}
		}
	}

	// make sure we still have a valid signing epoch
	i := len(f.network.NetworkEpochs) - 1
	now := time.Now().UTC()
	if now.After(f.network.NetworkEpochs[i].SignEnd) {
		return errors.New("netconf: no valid signing epoch found")
	}

	// now we make sure that in the present we have enough mint epochs (quorum)
	for ; i >= 0; i-- {
		if now.After(f.network.NetworkEpochs[i].SignStart) {
			break
		}
	}
	if i < 0 {
		i = 0
	}
	var q uint64
	for _, m := range f.mints {
		if len(m.MintEpochs) > i {
			q++
		}
	}
	if q < f.network.NetworkEpochs[i].QuorumM {
		return errors.New("netconf: not enough mint to reach quorum in present")
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
	f.network = n
	f.mints = make(map[string]*Mint)

	// TODO: the following is not correct in all situations, we have to validate
	// in the current signing period (and maybe later ones)

	for mn := range n.Mints() {
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
	if err := f.validate(); err != nil {
		return nil, err
	}
	return &f, nil
}
