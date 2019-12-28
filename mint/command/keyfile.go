package command

import (
	"github.com/frankbraun/codechain/command"
	"github.com/scritcash/scrit/util/homedir"
)

// KeyFile implements the scrit-mint 'keyfile' command.
func KeyFile(argv0 string, args ...string) error {
	return command.KeyFile("scrit", homedir.ScritMint(), argv0, args...)
}
