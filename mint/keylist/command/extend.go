package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/codechain/util/seckey"
	"github.com/scritcash/scrit/mint/identity"
	"github.com/scritcash/scrit/netconf"
	"github.com/scritcash/scrit/util/homedir"
)

func extend(net *netconf.Network, homeDir, secKey string) error {
	// load identity key
	sec, _, _, err := identity.Load(homeDir, secKey)
	if err != nil {
		return err
	}
	ik := netconf.NewIdentityKeyEd25519Priv(sec)

	id := ik.MarshalID()
	privFilename := filepath.Join(homeDir, netconf.DefPrivKeyListDir, id+".json")
	confFilename := filepath.Join(netconf.DefMintDir, id+".json")

	// make sure these files exist already and are valid
	mint, err := netconf.LoadMint(privFilename)
	if err != nil {
		return err
	}
	if _, err := netconf.LoadMint(confFilename); err != nil {
		return err
	}

	// extend key list
	if err := mint.Extend(ik, net); err != nil {
		return err
	}

	// save private key list
	if err := mint.Save(privFilename, 0700); err != nil {
		return err
	}
	// prune private keys
	mint.PrunePrivKeys()
	// save public configuration file
	return mint.Save(confFilename, 0755)
}

// Extend implements the scrit-mint 'keylist extend' command.
func Extend(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", argv0)
		fmt.Fprintf(os.Stderr, "Extend key list for %s.\n", netconf.DefNetConfFile)
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
	return extend(net, homeDir, *secKey)
}
