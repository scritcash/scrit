package command

import (
	"github.com/frankbraun/codechain/command"
	"github.com/scritcash/scrit/util/homedir"
)

// KeyGen implements the scrit-mint 'keygen' command.
func KeyGen(argv0 string, args ...string) error {
	return command.KeyGen("scrit", homedir.ScritMint(), argv0, args...)
}
