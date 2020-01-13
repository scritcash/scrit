package command

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/log"
	"github.com/scritcash/scrit/netconf"
)

const defaultPeriod = 30 * 24 * time.Hour

func start(
	filename string,
	m, n uint64,
	signStart, signEnd, validationEnd time.Time,
	mintIdentities []netconf.IdentityKey,
) error {
	exists, err := file.Exists(filename)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("file '%s' exists already", filename)
	}
	net := netconf.NewNetwork(m, n, signStart, signEnd, validationEnd,
		mintIdentities)
	if err := net.Validate(); err != nil {
		return err
	}
	return net.Save(filename)
}

// Start implements the scrit-gov 'start' command.
func Start(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s mint_identity [...]\n", argv0)
		fmt.Fprintf(os.Stderr, "Start new Scrit %s.\n", netconf.DefNetConfFile)
		fs.PrintDefaults()
	}
	m := fs.Uint64("m", 2, "The quorum m")
	n := fs.Uint64("n", 3, "Number of mints n")
	startSign := fs.String("start-sign", time.Now().UTC().Format(time.RFC3339),
		"Start of signing epoch")
	signingPeriod := fs.Duration("signing-period", defaultPeriod, "Signing period")
	validationPeriod := fs.Duration("validation-period", defaultPeriod, "Validation period")
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
	if fs.NArg() == 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	t, err := time.Parse(time.RFC3339, *startSign)
	if err != nil {
		return err
	}
	var keys []netconf.IdentityKey
	for _, arg := range fs.Args() {
		key, err := netconf.ParseIdentityKey(arg)
		if err != nil {
			return err
		}
		keys = append(keys, *key)
	}
	return start(netconf.DefNetConfFile, *m, *n, t, t.Add(*signingPeriod),
		t.Add(*signingPeriod).Add(*validationPeriod), keys)
}
