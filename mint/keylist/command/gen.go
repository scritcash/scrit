package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/codechain/util/seckey"
	"github.com/scritcash/scrit/mint/identity"
	"github.com/scritcash/scrit/netconf"
	"github.com/scritcash/scrit/util/homedir"
)

func gen(net *netconf.Network, homeDir, secKey string) error {
	sec, _, _, err := identity.Load(homeDir, secKey)
	if err != nil {
		return err
	}
	ik := netconf.NewIdentityKeyEd25519Priv(sec)
	filename := filepath.Join(netconf.DefMintDir, ik.MarshalID()+".json")
	exists, err := file.Exists(filename)
	if err != nil {
		return err
	}

	// TODO: continue here
	if exists {
		_, err := netconf.LoadMint(filename)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("file '%s' does not exist\n", filename)
	}

	return nil
}

// Gen implements the scrit-mint 'keylist gen' command.
func Gen(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
		fmt.Fprintf(os.Stderr, "Generate keylist for %s.\n", netconf.DefNetConfFile)
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
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	if err := secpkg.UpToDate("scrit"); err != nil {
		return err
	}
	homeDir := homedir.ScritMint()
	if err := seckey.Check(homeDir, *secKey); err != nil {
		return err
	}
	net, err := netconf.LoadNetwork(netconf.DefNetConfFile)
	if err != nil {
		return err
	}
	if err := net.Validate(); err != nil {
		return err
	}
	fmt.Println(net.Marshal())
	return gen(net, homeDir, *secKey)
}
