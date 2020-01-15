// scrit-engine is the low-level DBC engine for Scrit.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scritcash/scrit/engine/command"
)

func usage() {
	cmd := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage: %s validateconf [-d federation_dir]\n", cmd)
	os.Exit(2)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	argv0 := os.Args[0] + " " + os.Args[1]
	args := os.Args[2:]
	var err error
	switch os.Args[1] {
	case "validateconf":
		err = command.ValidateConf(argv0, args...)
	default:
		usage()
	}
	if err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
			os.Exit(1)
		}
		os.Exit(2)
	}
}
