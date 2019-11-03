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

// ErrSignEpochStartNotBeforeEnd is returned when the signing epoch start is not
// before the end.
var ErrSignEpochStartNotBeforeSignEnd = errors.New("netconf: signing epoch start is not before signing end")

// ErrSignEpochEndNotBeforeValidate is returned when the signing epoch start is
// not before the end.
var ErrSignEpochEndNotBeforeValidateEnd = errors.New("netconf: signing epoch end is not before end validate end")

// ErrSignEpochWrongBoundaires is returned when the signing epoch boundaries
// do not match exactly.
var ErrSignEpochWrongBoundaries = errors.New("netconf: signing epoch boundaries do not match exactly")
