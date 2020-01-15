package netconf

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/frankbraun/codechain/util/file"
	"github.com/scritcash/scrit/binencode"
)

// Mint defines the key list of a single mint for all epochs and where to
// reach the mint.
type Mint struct {
	Description     string       // description of mint (name)
	MintIdentityKey IdentityKey  // identity key of mint
	MintEpochs      []*MintEpoch // corresponding to global epochs
	URLs            []string     // how to reach the mint
}

// MintEpoch defines the key list of a single mint for a single epoch.
type MintEpoch struct {
	SignStart         time.Time     // start of signing epoch
	SignEnd           time.Time     // end of signing epoch
	ValidateEnd       time.Time     // end of validation epoch
	KeyList           []*SigningKey // the key list
	KeyListSignatures [][]byte      // signatures of key list (identity signature last)
}

func (m *Mint) generateKeys(ik *IdentityKey, n *Network, start int) error {
	dbcTypes := make(map[DBCType]bool)
	for i := start; i < len(n.NetworkEpochs); i++ {
		e := n.NetworkEpochs[i]
		for _, add := range e.DBCTypesAdded {
			dbcTypes[add] = true
		}
		for _, remove := range e.DBCTypesRemoved {
			delete(dbcTypes, remove)
		}
		dbcs := DBCTypeMapToSortedArray(dbcTypes)
		for _, dbc := range dbcs {
			pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				return err
			}
			sk := &SigningKey{
				Currency: dbc.Currency,
				Amount:   dbc.Amount,
				SigAlgo:  "ed25519", // TODO
				PubKey:   pubKey,
				PrivKey:  privKey,
			}
			m.MintEpochs[i].KeyList = append(m.MintEpochs[i].KeyList, sk)
		}
	}
	return m.sign(start)
}

func NewMint(
	description string,
	ik *IdentityKey,
	urls []string,
	n *Network,
) (*Mint, error) {
	var m Mint
	m.Description = description
	m.MintIdentityKey = *ik
	for _, ne := range n.NetworkEpochs {
		me := &MintEpoch{
			SignStart:   ne.SignStart,
			SignEnd:     ne.SignEnd,
			ValidateEnd: ne.ValidateEnd,
		}
		m.MintEpochs = append(m.MintEpochs, me)
	}
	m.URLs = urls
	if err := m.generateKeys(ik, n, 0); err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Mint) Extend(ik *IdentityKey, n *Network) error {
	start := 0
	for i, ne := range n.NetworkEpochs {
		if i < len(m.MintEpochs) {
			me := m.MintEpochs[i]
			if me.SignStart != ne.SignStart {
				return errors.New("netconf: epoch signature starts do not match")
			}
			if me.SignEnd != ne.SignEnd {
				return errors.New("netconf: epoch signature ends do not match")
			}
			if me.ValidateEnd != ne.ValidateEnd {
				return errors.New("netconf: epoch validation ends do not match")
			}
		} else {
			if start == 0 {
				start = i
			}
			me := &MintEpoch{
				SignStart:   ne.SignStart,
				SignEnd:     ne.SignEnd,
				ValidateEnd: ne.ValidateEnd,
			}
			m.MintEpochs = append(m.MintEpochs, me)
		}
	}
	return m.generateKeys(ik, n, start)
}

func (m *Mint) PrunePrivKeys() {
	for _, e := range m.MintEpochs {
		for _, k := range e.KeyList {
			k.PrivKey = nil
		}
	}
}

func (m *Mint) Save(filename string, perm os.FileMode) error {
	jsn, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err) // should never happen
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
	if err := ioutil.WriteFile(filename, jsn, perm); err != nil {
		return err
	}
	if exists {
		return os.Remove("." + filename)
	}
	return nil

}

// LoadMint loads a mint configuration from filename and return the
// Mint struct.
func LoadMint(filename string) (*Mint, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var mint Mint
	if err := json.Unmarshal(data, &mint); err != nil {
		return nil, err
	}
	return &mint, err
}

// sign all mint epochs in mint
func (m *Mint) sign(start int) error {
	for i := start; i < len(m.MintEpochs); i++ {
		e := m.MintEpochs[i]
		if err := e.sign(&m.MintIdentityKey); err != nil {
			return err
		}
	}
	return nil
}

// sign mint epoch.
func (me *MintEpoch) sign(ik *IdentityKey) error {
	encodingScheme := []interface{}{
		[]byte(ik.SigAlgo),
		ik.PubKey,
		me.SignStart.UTC().Unix(),
		me.SignEnd.UTC().Unix(),
		me.ValidateEnd.UTC().Unix(),
	}
	for _, k := range me.KeyList {
		encodingScheme = append(encodingScheme,
			[]byte(k.Currency),
			int64(k.Amount),
			[]byte(k.SigAlgo),
			k.PubKey,
		)
	}
	size, err := binencode.EncodeSize(encodingScheme...)
	if err != nil {
		return err
	}
	buf := make([]byte, size)
	enc, err := binencode.Encode(buf, encodingScheme...)
	if err != nil {
		return err
	}
	for _, k := range me.KeyList {
		sig := ed25519.Sign(k.PrivKey, enc)
		me.KeyListSignatures = append(me.KeyListSignatures, sig)
	}
	sig := ed25519.Sign(ik.privKey, enc)
	me.KeyListSignatures = append(me.KeyListSignatures, sig)
	return nil
}

// Verify mint epoch.
func (me *MintEpoch) Verify(ik *IdentityKey) error {
	// TODO
	return nil
}

// Validate the mint configuration.
func (mint *Mint) Validate() error {
	// TODO
	return nil
}
