package command

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/log"
	"github.com/scritcash/scrit/netconf"
)

func list(net *netconf.Network) error {
	var mints []string
	for m := range net.Mints() {
		mints = append(mints, m)
	}
	sort.Strings(mints)
	for _, m := range mints {
		fmt.Println(m)
	}
	return nil
}

// List implements the scrit-gov 'mint list' command.
func List(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
		fmt.Fprintf(os.Stderr, "List all mints that will be active in the last epoch of %s.\n",
			netconf.DefNetConfFile)
		fs.PrintDefaults()
	}
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	if err := secpkg.UpToDate("scrit"); err != nil {
		return err
	}
	// load
	net, err := netconf.LoadNetwork(netconf.DefNetConfFile)
	if err != nil {
		return err
	}
	// validate
	if err := net.Validate(); err != nil {
		return err
	}
	// list
	if err := list(net); err != nil {
		return err
	}
	return nil
}
