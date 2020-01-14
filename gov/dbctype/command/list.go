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
	var dbcTypes []netconf.DBCType
	for t := range net.DBCTypes() {
		dbcTypes = append(dbcTypes, t)
	}
	sort.Slice(dbcTypes, func(i, j int) bool {
		if dbcTypes[i].Currency < dbcTypes[j].Currency ||
			(dbcTypes[i].Currency == dbcTypes[j].Currency &&
				dbcTypes[i].Amount < dbcTypes[j].Amount) {
			return true
		}
		return false
	})
	for _, t := range dbcTypes {
		fmt.Printf("%s\t%d\n", t.Currency, t.Amount)
	}
	return nil
}

// List implements the scrit-gov 'dbctype list' command.
func List(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
		fmt.Fprintf(os.Stderr, "List all DBC types that will be active in the last epoch of %s.\n",
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
