package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util"
	"github.com/frankbraun/codechain/util/log"
	"github.com/scritcash/scrit/netconf"
)

// Status implements the scrit-gov 'status' command.
func Status(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s]\n", argv0)
		fmt.Fprintf(os.Stderr, "Print status of %s.\n", netconf.DefNetConfFile)
		fs.PrintDefaults()
	}
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	if err := secpkg.UpToDate("scrit"); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	net, err := netconf.LoadNetwork(netconf.DefNetConfFile)
	if err != nil {
		util.Fatal(err)
	}
	if err := net.Validate(); err != nil {
		util.Fatal(err)
	}
	fmt.Println(net.Marshal())
	return nil
}
