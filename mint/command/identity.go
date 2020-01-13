package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/def"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/codechain/util/seckey"
	"github.com/scritcash/scrit/netconf"
	"github.com/scritcash/scrit/util/homedir"
)

func identity(homeDir, secKey string) error {
	if secKey == "" {
		secretDir := filepath.Join(homeDir, def.SecretsSubDir)
		files, err := ioutil.ReadDir(secretDir)
		if err != nil {
			return err
		}
		if len(files) > 1 {
			return fmt.Errorf("directory '%s' contains more than one secret file, use option -s",
				secretDir)
		}
		secKey = filepath.Join(secretDir, files[0].Name())
	}
	sec, _, comment, err := seckey.Read(secKey)
	if err != nil {
		return err
	}
	ik := netconf.NewIdentityKeyEd25519Priv(sec)
	fmt.Println(ik.MarshalID())
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
	if err := seckey.Check(homedir.ScritMint(), *secKey); err != nil {
		return err
	}
	if err := secpkg.UpToDate("scrit"); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	return identity(homeDir, *secKey)
}
