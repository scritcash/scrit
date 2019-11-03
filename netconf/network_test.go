package netconf

import (
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	net, err := Load(filepath.Join("testdata", DefNetConfFile))
	if err != nil {
		t.Fatal(err)
	}
	if err := net.Validate(); err != nil {
		t.Fatal(err)
	}
}
