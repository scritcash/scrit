// Package identity defines helper functions related to mint identity keys.
package identity

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/frankbraun/codechain/util/def"
	"github.com/frankbraun/codechain/util/seckey"
)

// Load secret key from homeDir/def.SecretsSubdir.
// If secKey is empty a secret key from homeDir/def.SecretsSubdir is loaded
// only if it contains exactly one secret.
func Load(homeDir, secKey string) (*[64]byte, *[64]byte, []byte, error) {
	if secKey == "" {
		secretDir := filepath.Join(homeDir, def.SecretsSubDir)
		files, err := ioutil.ReadDir(secretDir)
		if err != nil {
			return nil, nil, nil, err
		}
		if len(files) > 1 {
			return nil, nil, nil,
				fmt.Errorf("directory '%s' contains more than one secret file, use option -s",
					secretDir)
		}
		secKey = filepath.Join(secretDir, files[0].Name())
	}
	return seckey.Read(secKey)
}
