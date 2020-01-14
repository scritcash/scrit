package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/log"
	"github.com/scritcash/scrit/netconf"
)

func setQuorum(net *netconf.Network, m uint64) error {
	// make sure network has a future epoch
	if err := net.HasFuture(); err != nil {
		return err
	}
	net.SetQuorum(m)
	return nil
}

// SetQuorum implements the scrit-gov 'epoch setquorum' command.
func SetQuorum(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
		fmt.Fprintf(os.Stderr, "Set quorum for future epoch of %s.\n", netconf.DefNetConfFile)
		fs.PrintDefaults()
	}
	m := fs.Uint64("m", 0, "The quorum m")
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	if *m == 0 {
		fmt.Fprintf(os.Stderr, "%s: option -m is mandatory\n", argv0)
		return flag.ErrHelp
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
	// edit
	if err := setQuorum(net, *m); err != nil {
		return err
	}
	// validate again
	if err := net.Validate(); err != nil {
		return err
	}
	// save
	if err := net.Save(netconf.DefNetConfFile); err != nil {
		return err
	}
	return nil
}
