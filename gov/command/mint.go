package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/scritcash/scrit/gov/mint/command"
)

func usageMint(cmd string) error {
	fmt.Fprintf(os.Stderr, "Usage: %s add\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s remove\n", cmd)
	return flag.ErrHelp
}

// Mint implements the scrit-gov 'mint' command.
func Mint(argv0 string, args ...string) error {
	if len(args) < 1 {
		return usageMint(argv0)
	}
	newArgv0 := argv0 + " " + args[0]
	newArgs := args[1:]
	switch args[0] {
	case "add":
		return command.Add(newArgv0, newArgs...)
	case "remove":
		return command.Remove(newArgv0, newArgs...)
	default:
		return usageMint(argv0)
	}
}
