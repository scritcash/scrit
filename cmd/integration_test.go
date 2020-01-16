package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/frankbraun/codechain/command"
	"github.com/frankbraun/codechain/util/seckey"
	scritEngine "github.com/scritcash/scrit/engine/command"
	scritGov "github.com/scritcash/scrit/gov/command"
	scritDBCType "github.com/scritcash/scrit/gov/dbctype/command"
	scritEpoch "github.com/scritcash/scrit/gov/epoch/command"
	scritMint "github.com/scritcash/scrit/mint/command"
	scritKeyList "github.com/scritcash/scrit/mint/keylist/command"
)

// Test setting up a federation of Scrit mints (see doc/federation-setup.md).
func TestFederationSetup(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "scrit_integration_test")
	if err != nil {
		t.Fatalf("ioutil.TempDir() failed: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	// create separate mint directories
	mint1dir := filepath.Join(tmpdir, "mint1")
	mint2dir := filepath.Join(tmpdir, "mint2")
	mint3dir := filepath.Join(tmpdir, "mint3")
	if err := os.Mkdir(mint1dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(mint2dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(mint3dir, 0755); err != nil {
		t.Fatal(err)
	}

	// create identity key for mint 1
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint1dir); err != nil {
		t.Fatal(err)
	}
	seckey.TestPass = "test"
	command.TestComment = "Test Mint 1"
	if err := scritMint.KeyGen("scrit-mint keygen"); err != nil {
		t.Fatal(err)
	}

	// create identity key for mint 2
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint2dir); err != nil {
		t.Fatal(err)
	}
	command.TestComment = "Test Mint 2"
	if err := scritMint.KeyGen("scrit-mint keygen"); err != nil {
		t.Fatal(err)
	}

	// create identity key for mint 3
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint3dir); err != nil {
		t.Fatal(err)
	}
	command.TestComment = "Test Mint 3"
	if err := scritMint.KeyGen("scrit-mint keygen"); err != nil {
		t.Fatal(err)
	}

	// get identity key for mint 1
	stdout := os.Stdout
	tmpfile1, err := ioutil.TempFile("", "scrit_integration_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile1.Name())
	os.Stdout = tmpfile1
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint1dir); err != nil {
		t.Fatal(err)
	}
	if err := scritMint.Identity("scrit-mint identity"); err != nil {
		t.Fatal(err)
	}
	os.Stdout = stdout
	if err := tmpfile1.Close(); err != nil {
		t.Fatal(err)
	}
	buf, err := ioutil.ReadFile(tmpfile1.Name())
	if err != nil {
		t.Fatal(err)
	}
	lines := bytes.Split(buf, []byte("\n"))
	key1 := string(lines[len(lines)-2])

	// get identity key for mint 2
	tmpfile2, err := ioutil.TempFile("", "scrit_integration_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())
	os.Stdout = tmpfile2
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint2dir); err != nil {
		t.Fatal(err)
	}
	if err := scritMint.Identity("scrit-mint identity"); err != nil {
		t.Fatal(err)
	}
	os.Stdout = stdout
	if err := tmpfile2.Close(); err != nil {
		t.Fatal(err)
	}
	buf, err = ioutil.ReadFile(tmpfile2.Name())
	if err != nil {
		t.Fatal(err)
	}
	lines = bytes.Split(buf, []byte("\n"))
	key2 := string(lines[len(lines)-2])

	// get identity key for mint 3
	tmpfile3, err := ioutil.TempFile("", "scrit_integration_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile3.Name())
	os.Stdout = tmpfile3
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint3dir); err != nil {
		t.Fatal(err)
	}
	if err := scritMint.Identity("scrit-mint identity"); err != nil {
		t.Fatal(err)
	}
	os.Stdout = stdout
	if err := tmpfile3.Close(); err != nil {
		t.Fatal(err)
	}
	buf, err = ioutil.ReadFile(tmpfile3.Name())
	if err != nil {
		t.Fatal(err)
	}
	lines = bytes.Split(buf, []byte("\n"))
	key3 := string(lines[len(lines)-2])

	// setup the federation (2-of-3):
	if err := os.Chdir(tmpdir); err != nil {
		t.Fatal(err)
	}
	err = scritGov.Start("scrit-gov start", "-m", "2", "-n", "3", key1, key2, key3)
	if err != nil {
		t.Fatal(err)
	}

	// define first DBC types (in denominations of 1, 2, and 5 EUR)
	err = scritDBCType.Add("scrit-gov dbctype add", "-currency", "EUR", "-amount", "100000000")
	if err != nil {
		t.Fatal(err)
	}
	err = scritDBCType.Add("scrit-gov dbctype add", "-currency", "EUR", "-amount", "200000000")
	if err != nil {
		t.Fatal(err)
	}
	err = scritDBCType.Add("scrit-gov dbctype add", "-currency", "EUR", "-amount", "500000000")
	if err != nil {
		t.Fatal(err)
	}

	// create key lists
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint1dir); err != nil {
		t.Fatal(err)
	}
	if err := scritKeyList.Create("scrit-mint keylist create", "-desc", "mint1", "https://mint1.example.com"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint2dir); err != nil {
		t.Fatal(err)
	}
	if err := scritKeyList.Create("scrit-mint keylist create", "-desc", "mint2", "https://mint2.example.net"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint3dir); err != nil {
		t.Fatal(err)
	}
	if err := scritKeyList.Create("scrit-mint keylist create", "-desc", "mint3", "https://mint3.example.org"); err != nil {
		t.Fatal(err)
	}

	// define the second signing epoch
	err = scritEpoch.Add("scrit-gov epoch add")
	if err != nil {
		t.Fatal(err)
	}

	// extend key lists
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint1dir); err != nil {
		t.Fatal(err)
	}
	if err := scritKeyList.Extend("scrit-mint keylist extend"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint2dir); err != nil {
		t.Fatal(err)
	}
	if err := scritKeyList.Extend("scrit-mint keylist extend"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint3dir); err != nil {
		t.Fatal(err)
	}
	if err := scritKeyList.Extend("scrit-mint keylist extend"); err != nil {
		t.Fatal(err)
	}

	// test configuration
	if err := scritEngine.ValidateConf("scrit-engine validateconf"); err != nil {
		t.Fatal(err)
	}

	if err := scritDBCType.List("scrit-gov dbctype list"); err != nil {
		t.Fatal(err)
	}
}
