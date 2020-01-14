// Package def defines default values used in Scrit.
package def

import (
	"time"
)

// SigningPeriod defines the default length of a signing epoch.
const SigningPeriod = 30 * 24 * time.Hour // 30d

// ValidationPeriod defines the default length of a validation epoch.
const ValidationPeriod = 30 * 24 * time.Hour // 30d
