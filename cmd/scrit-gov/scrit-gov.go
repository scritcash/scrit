/// scrit-mint is a Scrit mint.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scritcash/scrit/gov/command"
)

func usage() {
	cmd := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage: %s start\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s dbctype\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s epoch\n", cmd)
	fmt.Fprintf(os.Stderr, "       %s status\n", cmd)
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
	case "start":
		err = command.Start(argv0, args...)
	case "dbctype":
		err = command.DBCType(argv0, args...)
	case "epoch":
		err = command.Epoch(argv0, args...)
	case "status":
		err = command.Status(argv0, args...)
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
