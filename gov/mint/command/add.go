package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/log"
	"github.com/scritcash/scrit/netconf"
)

func add(net *netconf.Network, key *netconf.IdentityKey) error {
	// make sure network has a future epoch
	if err := net.HasFuture(); err != nil {
		return err
	}
	// make sure mint has not been added yet
	mints := net.Mints()
	if mints[key.MarshalID()] {
		return fmt.Errorf("mint already added: %v", key.MarshalID())
	}
	// add mint identity key
	net.MintAdd(key)
	return nil
}

// Add implements the scrit-gov 'epoch add' command.
func Add(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s mint_identity\n", argv0)
		fmt.Fprintf(os.Stderr, "Add new mint_identity to %s.\n", netconf.DefNetConfFile)
		fs.PrintDefaults()
	}
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
	// parse key
	key, err := netconf.ParseIdentityKey(fs.Arg(0))
	if err != nil {
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
	if err := add(net, key); err != nil {
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
