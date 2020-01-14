package command

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/log"
	"github.com/scritcash/scrit/netconf"
	"github.com/scritcash/scrit/util/def"
)

func add(net *netconf.Network, siginingPeriod, validationPeriod time.Duration) error {
	// TODO: check if signingPeriod and/or validationPeriod changes from the
	// period used in the last epoch.
	net.EpochAdd(siginingPeriod, validationPeriod)
	return nil
}

// Add implements the scrit-gov 'epoch add' command.
func Add(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
		fmt.Fprintf(os.Stderr, "Add new DBC type to future epoch of %s.\n", netconf.DefNetConfFile)
		fs.PrintDefaults()
	}
	signingPeriod := fs.Duration("signing-period", def.SigningPeriod, "Length of signing period")
	validationPeriod := fs.Duration("validation-period", def.ValidationPeriod, "Length of validation period")
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
	// edit
	if err := add(net, *signingPeriod, *validationPeriod); err != nil {
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
