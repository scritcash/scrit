package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/codechain/util/seckey"
	"github.com/scritcash/scrit/mint/identity"
	"github.com/scritcash/scrit/netconf"
	"github.com/scritcash/scrit/util/homedir"
)

func showIdentity(homeDir, secKey string) error {
	sec, _, comment, err := identity.Load(homeDir, secKey)
	if err != nil {
		return err
	}
	ik := netconf.NewIdentityKeyEd25519Priv(sec)
	fmt.Println(string(comment))
	fmt.Println(ik.MarshalID()) // this must be the last output line!
	return nil
}

// Identity implements the scrit-mint 'identity' command.
func Identity(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-s seckey.bin]\n", argv0)
		fmt.Fprintf(os.Stderr, "Print mint identity.\n")
		fs.PrintDefaults()
	}
	secKey := fs.String("s", "", "Secret key file")
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	homeDir := homedir.ScritMint()
	if err := seckey.Check(homeDir, *secKey); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	if err := secpkg.UpToDate("scrit"); err != nil {
		return err
	}
	return showIdentity(homeDir, *secKey)
}
