package netconf

import (
	"errors"
)

// ErrZeroM is returned when the quorum M is 0.
var ErrZeroM = errors.New("netconf: quorum M is 0")

// ErrZeroN is returned when the number of mints N is 0.
var ErrZeroN = errors.New("netconf: number of mints N is 0")

// ErrMGreaterN is returned when the quorum M is greater than the number of
// mints N.
var ErrMGreaterN = errors.New("netconf: quorum M is greater than the number of mints N")

// ErrQuorumTooSmall is returned when the qurorum is smaller or equal the
// number of mints N divided by 2.
var ErrQuorumTooSmall = errors.New("netconf: quorum M too small (must be > n/2)")

// ErrSignEpochStartNotBeforeSignEnd is returned when the signing epoch start is not
// before the end.
var ErrSignEpochStartNotBeforeSignEnd = errors.New("netconf: signing epoch start is not before signing end")

// ErrSignEpochEndNotBeforeValidateEnd is returned when the signing epoch start is
// not before the end.
var ErrSignEpochEndNotBeforeValidateEnd = errors.New("netconf: signing epoch end is not before end validate end")

// ErrSignEpochWrongBoundaries is returned when the signing epoch boundaries
// do not match exactly.
var ErrSignEpochWrongBoundaries = errors.New("netconf: signing epoch boundaries do not match exactly")

// ErrValidationLongerThanNextSigning is returned when the validation period is
// longer than the next signing epoch.
var ErrValidationLongerThanNextSigning = errors.New("netconf: validation period is longer than next signing epoch")

// ErrNoFuture is returned when a network has no epoch which starts in the future.
var ErrNoFuture = errors.New("netconf: network has no epoch which starts in the future")

// ErrDBCTypesOverlap is returned if the sets DBCTypesAdded and DBCTypesRemoved overlap.
var ErrDBCTypesOverlap = errors.New("netconf: DBCTypesAdded and DBCTypesRemoved overlap.")
