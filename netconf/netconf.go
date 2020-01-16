// Package netconf implements the Scrit network configuration.
package netconf

import (
	"time"
)

// DefNetConfFile defines the default network configuration filename.
const DefNetConfFile = "federation.json"

// DefMintDir defines the default sub-directory for mint configurations.
const DefMintDir = "mints"

// DefDBCDir defines the default sub-directory for DBC creation and
// destruction lists.
const DefDBCDir = "dbcs"

// DefPrivKeyListDir defines the default sub-directory for private mint key lists.
const DefPrivKeyListDir = "privkeylists"

// DefDBCCreate defines the name of the list of DBCs to be created.
const DefDBCCreate = "create.json"

// DefDBCDestroyed defines the name of the list of destroyed DBCs.
const DefDBCDestroyed = "destroyed.json"

// DefStartTime defines the default signing start: tomorrow at midnight.
func DefStartTime() time.Time {
	now := time.Now().UTC()
	year, month, day := now.Date()
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	t = t.Add(time.Hour * 48)
	return t
}
