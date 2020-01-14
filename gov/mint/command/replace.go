package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/log"
	"github.com/scritcash/scrit/netconf"
)

func replace(
	net *netconf.Network,
	newKey, oldKey *netconf.IdentityKey,
	sig string,
) error {
	// make sure network has a future epoch
	if err := net.HasFuture(); err != nil {
		return err
	}
	// make sure mint has been added before
	mints := net.Mints()
	if mints[newKey.MarshalID()] {
		return fmt.Errorf("mint to replace to already added: %v", newKey.MarshalID())
	}
	if !mints[oldKey.MarshalID()] {
		return fmt.Errorf("mint to replace from not added before: %v", oldKey.MarshalID())
	}
	// remove mint identity key
	net.MintReplace(newKey, oldKey, sig)
	return nil
}

// Replace implements the scrit-gov 'mint replace' command.
func Replace(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s new_mint_id old_mint_id signature\n", argv0)
		fmt.Fprintf(os.Stderr, "Replace old_mint_id with new_mint_id in %s.\n", netconf.DefNetConfFile)
		fmt.Fprintf(os.Stderr, "The signature is by old_mint_id over new_mint_id.\n")
		fs.PrintDefaults()
	}
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	if fs.NArg() != 3 {
		fs.Usage()
		return flag.ErrHelp
	}
	if err := secpkg.UpToDate("scrit"); err != nil {
		return err
	}
	// parse keys
	newKey, err := netconf.ParseIdentityKey(fs.Arg(0))
	if err != nil {
		return err
	}
	oldKey, err := netconf.ParseIdentityKey(fs.Arg(1))
	if err != nil {
		return err
	}
	sig := fs.Arg(2)
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
	if err := replace(net, newKey, oldKey, sig); err != nil {
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
