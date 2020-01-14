package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/scritcash/scrit/gov/dbctype/command"
)

func usage(cmd string) error {
	fmt.Fprintf(os.Stderr, "Usage: %s add\n", cmd)
	return flag.ErrHelp
}

// DBCType implements the scrit-gov 'dbctype' command.
func DBCType(argv0 string, args ...string) error {
	if len(args) < 1 {
		return usage(argv0)
	}
	newArgv0 := argv0 + " " + args[0]
	newArgs := args[1:]
	switch args[0] {
	case "add":
		return command.Add(newArgv0, newArgs...)
	default:
		return usage(argv0)
	}
}
