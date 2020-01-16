package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/log"
	"github.com/scritcash/scrit/netconf"
)

func reissue(fed *netconf.Federation, dbc string) error {
	// TODO
	return nil
}

// Reissue implements the scrit-engine 'reissue' command.
func Reissue(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-d federation_dir] DBC\n", argv0)
		fmt.Fprintf(os.Stderr, "Reissue DBC.\n")
		fs.PrintDefaults()
	}
	dir := fs.String("d", ".", "Set federation directory")
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return flag.ErrHelp
	}
	if err := secpkg.UpToDate("scrit"); err != nil {
		return err
	}
	fed, err := netconf.LoadFederation(*dir)
	if err != nil {
		return err
	}
	dbc := fs.Arg(0)
	return reissue(fed, dbc)
}
