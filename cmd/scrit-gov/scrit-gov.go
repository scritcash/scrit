// scrit-gov is a government helper tool for Scrit.
package main

import (
	"fmt"

	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/util"
)

func main() {
	if err := secpkg.UpToDate("scrit"); err != nil {
		util.Fatal(err)
	}
	fmt.Println("scrit-gov")
}
