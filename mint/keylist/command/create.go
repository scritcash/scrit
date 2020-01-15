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

func create(
	net *netconf.Network,
	homeDir, secKey, desc string,
	urls []string,
) error {
	// load identity key
	sec, _, _, err := identity.Load(homeDir, secKey)
	if err != nil {
		return err
	}
	ik := netconf.NewIdentityKeyEd25519Priv(sec)

	// make sure the '~/.config/scrit-mint/privkeylists' directory exists
	if err := os.MkdirAll(filepath.Join(homeDir, netconf.DefPrivKeyListDir), 0755); err != nil {
		return err
	}
	// make sure the 'mints' directory exists
	if err := os.MkdirAll(netconf.DefMintDir, 0755); err != nil {
		return err
	}

	id := ik.MarshalID()
	privFilename := filepath.Join(homeDir, netconf.DefPrivKeyListDir, id+".json")
	confFilename := filepath.Join(netconf.DefMintDir, id+".json")

	// make sure these files do not exist already
	exists, err := file.Exists(privFilename)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("file '%s' exists already", privFilename)
	}
	exists, err = file.Exists(confFilename)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("file '%s' exists already", confFilename)
	}

	// create key list and configuration file
	mint, err := netconf.NewMint(desc, ik, urls, net)
	if err != nil {
		return err
	}
	// save private key list
	if err := mint.Save(privFilename); err != nil {
		return err
	}
	// prune private keys
	mint.PrunePrivKeys()
	// save public configuration file
	return mint.Save(confFilename)
}

// Create implements the scrit-mint 'keylist create' command.
func Create(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s url [...]\n", argv0)
		fmt.Fprintf(os.Stderr, "Create keylist for %s.\n", netconf.DefNetConfFile)
		fs.PrintDefaults()
	}
	desc := fs.String("desc", "", "Description of mint (name)")
	secKey := fs.String("s", "", "Secret key file")
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	if *desc == "" {
		return fmt.Errorf("%s: option -desc is mandatory", argv0)
	}
	if fs.NArg() == 0 {
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
	return create(net, homeDir, *secKey, *desc, fs.Args())
}
