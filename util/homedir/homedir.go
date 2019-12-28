// Package homedir implements helper methods to get the home directories of
// various tools.
package homedir

import (
	"github.com/frankbraun/codechain/util/homedir"
)

// ScritMint returns the home directory for 'scrit-mint'.
func ScritMint() string {
	return homedir.Get("scrit-mint")
}
