package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/scritcash/scrit/mint/keylist/command"
)

func usageKeyList(cmd string) error {
	fmt.Fprintf(os.Stderr, "Usage: %s gen\n", cmd)
	return flag.ErrHelp
}

// KeyList implements the scrit-mint 'keylist' command.
func KeyList(argv0 string, args ...string) error {
	if len(args) < 1 {
		return usageKeyList(argv0)
	}
	newArgv0 := argv0 + " " + args[0]
	newArgs := args[1:]
	switch args[0] {
	case "gen":
		return command.Gen(newArgv0, newArgs...)
	default:
		return usageKeyList(argv0)
	}
}
