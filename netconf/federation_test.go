package netconf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/scritcash/scrit/util/def"
)

func TestLoadFederation(t *testing.T) {
	_, err := LoadFederation("testdata")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFederation(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "scrit_federation_test")
	if err != nil {
		t.Fatalf("ioutil.TempDir() failed: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	ik1, err := NewIdentityKey()
	if err != nil {
		t.Fatal(err)
	}
	ik2, err := NewIdentityKey()
	if err != nil {
		t.Fatal(err)
	}
	ik3, err := NewIdentityKey()
	if err != nil {
		t.Fatal(err)
	}
	iks := []IdentityKey{*ik1, *ik2, *ik3}
	start := DefStartTime()
	net := NewNetwork(2, 3, start, start.Add(def.SigningPeriod),
		start.Add(def.SigningPeriod).Add(def.ValidationPeriod), iks)

	m1, err := NewMint("mint1", ik1, []string{"https:\\mint1.example.com"}, net)
	if err != nil {
		t.Fatal(err)
	}
	m2, err := NewMint("mint2", ik1, []string{"https:\\mint2.example.net"}, net)
	if err != nil {
		t.Fatal(err)
	}
	m3, err := NewMint("mint3", ik1, []string{"https:\\mint3.example.org"}, net)
	if err != nil {
		t.Fatal(err)
	}

	if err := net.Validate(); err != nil {
		t.Error(err)
	}
	if err := m1.Validate(); err != nil {
		t.Error(err)
	}
	if err := m2.Validate(); err != nil {
		t.Error(err)
	}
	if err := m3.Validate(); err != nil {
		t.Error(err)
	}

	net.EpochAdd(def.SigningPeriod, def.ValidationPeriod)
	if err := m1.Extend(ik1, net); err != nil {
		t.Error(err)
	}
	if err := m2.Extend(ik2, net); err != nil {
		t.Error(err)
	}
	if err := m3.Extend(ik3, net); err != nil {
		t.Error(err)
	}

	if err := net.Validate(); err != nil {
		t.Error(err)
	}
	if err := m1.Validate(); err != nil {
		t.Error(err)
	}
	if err := m2.Validate(); err != nil {
		t.Error(err)
	}
	if err := m3.Validate(); err != nil {
		t.Error(err)
	}

	m1.PrunePrivKeys()
	m2.PrunePrivKeys()
	m3.PrunePrivKeys()

	if err := m1.Validate(); err != nil {
		t.Error(err)
	}
	if err := m2.Validate(); err != nil {
		t.Error(err)
	}
	if err := m3.Validate(); err != nil {
		t.Error(err)
	}

	filename := filepath.Join(tmpdir, DefNetConfFile)
	if err := net.Save(filename); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(tmpdir, DefMintDir), 0755); err != nil {
		t.Fatal(err)
	}
	filename = filepath.Join(tmpdir, DefMintDir, ik1.MarshalID()+".json")
	if err := m1.Save(filename, 0755); err != nil {
		t.Fatal(err)
	}
	filename = filepath.Join(tmpdir, DefMintDir, ik2.MarshalID()+".json")
	if err := m2.Save(filename, 0755); err != nil {
		t.Fatal(err)
	}
	filename = filepath.Join(tmpdir, DefMintDir, ik3.MarshalID()+".json")
	if err := m3.Save(filename, 0755); err != nil {
		t.Fatal(err)
	}
	_, err = LoadFederation(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
}
