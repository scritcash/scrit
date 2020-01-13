package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/frankbraun/codechain/command"
	"github.com/frankbraun/codechain/util/seckey"
	scritMint "github.com/scritcash/scrit/mint/command"
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
		t.Error(err)
	}
	if err := os.Mkdir(mint2dir, 0755); err != nil {
		t.Error(err)
	}
	if err := os.Mkdir(mint3dir, 0755); err != nil {
		t.Error(err)
	}

	// create identity key for mint 1
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint1dir); err != nil {
		t.Error(err)
	}
	seckey.TestPass = "test"
	command.TestComment = "Test Mint 1"
	if err := scritMint.KeyGen("scrit-mint keygen"); err != nil {
		t.Error(err)
	}

	// create identity key for mint 2
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint2dir); err != nil {
		t.Error(err)
	}
	command.TestComment = "Test Mint 2"
	if err := scritMint.KeyGen("scrit-mint keygen"); err != nil {
		t.Error(err)
	}

	// create identity key for mint 3
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint3dir); err != nil {
		t.Error(err)
	}
	command.TestComment = "Test Mint 3"
	if err := scritMint.KeyGen("scrit-mint keygen"); err != nil {
		t.Error(err)
	}

	// get identity key for mint 1
	stdout := os.Stdout
	tmpfile1, err := ioutil.TempFile("", "scrit_integration_test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile1.Name())
	os.Stdout = tmpfile1
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint1dir); err != nil {
		t.Error(err)
	}
	if err := scritMint.Identity("scrit-mint identity"); err != nil {
		t.Error(err)
	}
	os.Stdout = stdout
	if err := tmpfile1.Close(); err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadFile(tmpfile1.Name())
	if err != nil {
		t.Error(err)
	}
	lines := bytes.Split(buf, []byte("\n"))
	key1 := string(lines[len(lines)-2])

	// get identity key for mint 2
	tmpfile2, err := ioutil.TempFile("", "scrit_integration_test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile2.Name())
	os.Stdout = tmpfile2
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint2dir); err != nil {
		t.Error(err)
	}
	if err := scritMint.Identity("scrit-mint identity"); err != nil {
		t.Error(err)
	}
	os.Stdout = stdout
	if err := tmpfile2.Close(); err != nil {
		t.Error(err)
	}
	buf, err = ioutil.ReadFile(tmpfile2.Name())
	if err != nil {
		t.Error(err)
	}
	lines = bytes.Split(buf, []byte("\n"))
	key2 := string(lines[len(lines)-2])

	// get identity key for mint 3
	tmpfile3, err := ioutil.TempFile("", "scrit_integration_test")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpfile3.Name())
	os.Stdout = tmpfile3
	if err := os.Setenv("SCRIT-MINTHOMEDIR", mint3dir); err != nil {
		t.Error(err)
	}
	if err := scritMint.Identity("scrit-mint identity"); err != nil {
		t.Error(err)
	}
	os.Stdout = stdout
	if err := tmpfile3.Close(); err != nil {
		t.Error(err)
	}
	buf, err = ioutil.ReadFile(tmpfile3.Name())
	if err != nil {
		t.Error(err)
	}
	lines = bytes.Split(buf, []byte("\n"))
	key3 := string(lines[len(lines)-2])

	fmt.Println(key1)
	fmt.Println(key2)
	fmt.Println(key3)
}