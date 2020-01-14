package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/scritcash/scrit/gov/epoch/command"
)

func usageEpoch(cmd string) error {
	fmt.Fprintf(os.Stderr, "Usage: %s add\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s setquorum\n", cmd)
	return flag.ErrHelp
}

// Epoch implements the scrit-gov 'epoch' command.
func Epoch(argv0 string, args ...string) error {
	if len(args) < 1 {
		return usageEpoch(argv0)
	}
	newArgv0 := argv0 + " " + args[0]
	newArgs := args[1:]
	switch args[0] {
	case "add":
		return command.Add(newArgv0, newArgs...)
	case "setquorum":
		return command.SetQuorum(newArgv0, newArgs...)
	default:
		return usageEpoch(argv0)
	}
}
