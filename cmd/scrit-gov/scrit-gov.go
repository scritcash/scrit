// scrit-gov is a government helper tool for Scrit.
package main

import (
	"fmt"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util"
	"github.com/scritcash/scrit/netconf"
)

func main() {
	if err := secpkg.UpToDate("scrit"); err != nil {
		util.Fatal(err)
	}
	net, err := netconf.Load(netconf.DefNetConfFile)
	if err != nil {
		util.Fatal(err)
	}
	if err := net.Validate(); err != nil {
		util.Fatal(err)
	}
	fmt.Println(net.Marshal())
}
