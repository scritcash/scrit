package netconf

import (
	"testing"
	"time"
)

var (
	t1 time.Time
	t2 time.Time
	t3 time.Time
	t4 time.Time
)

func init() {
	var err error
	t1, err = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		panic(err)
	}
	t2, err = time.Parse(time.RFC3339, "2006-02-02T15:04:05Z")
	if err != nil {
		panic(err)
	}
	t3, err = time.Parse(time.RFC3339, "2006-03-02T15:04:05Z")
	if err != nil {
		panic(err)
	}
	t4, err = time.Parse(time.RFC3339, "2006-04-02T15:04:05Z")
	if err != nil {
		panic(err)
	}
}

func TestEpochValidate(t *testing.T) {
	testCases := []struct {
		epoch     NetworkEpoch
		errorCode error
	}{
		{
			NetworkEpoch{
				QuorumM:        0,
				NumberOfMintsN: 1,
				SignStart:      t1,
				SignEnd:        t2,
				ValidateEnd:    t3,
			},
			ErrZeroM,
		},
		{
			NetworkEpoch{
				QuorumM:        1,
				NumberOfMintsN: 0,
				SignStart:      t1,
				SignEnd:        t2,
				ValidateEnd:    t3,
			},
			ErrZeroN,
		},
		{
			NetworkEpoch{
				QuorumM:        2,
				NumberOfMintsN: 1,
				SignStart:      t1,
				SignEnd:        t2,
				ValidateEnd:    t3,
			},
			ErrMGreaterN,
		},
		{
			NetworkEpoch{
				QuorumM:        5,
				NumberOfMintsN: 10,
				SignStart:      t1,
				SignEnd:        t2,
				ValidateEnd:    t3,
			},
			ErrQuorumTooSmall,
		},
		{
			NetworkEpoch{
				QuorumM:        6,
				NumberOfMintsN: 10,
				SignStart:      t1,
				SignEnd:        t2,
				ValidateEnd:    t3,
			},
			nil,
		},
		{
			NetworkEpoch{
				QuorumM:        5,
				NumberOfMintsN: 11,
				SignStart:      t1,
				SignEnd:        t2,
				ValidateEnd:    t3,
			},
			ErrQuorumTooSmall,
		},
		{
			NetworkEpoch{
				QuorumM:        6,
				NumberOfMintsN: 11,
				SignStart:      t1,
				SignEnd:        t2,
				ValidateEnd:    t3,
			},
			nil,
		},
		{
			NetworkEpoch{
				QuorumM:        8,
				NumberOfMintsN: 10,
				SignStart:      t1,
				SignEnd:        t1,
			},
			ErrSignEpochStartNotBeforeSignEnd,
		},
		{
			NetworkEpoch{
				QuorumM:        8,
				NumberOfMintsN: 10,
				SignStart:      t1,
				SignEnd:        t2,
				ValidateEnd:    t2,
			},
			ErrSignEpochEndNotBeforeValidateEnd,
		},
	}
	for _, testCase := range testCases {
		err := testCase.epoch.Validate()
		if err != testCase.errorCode {
			if err != testCase.errorCode {
				t.Fatalf("Validate(%#v) should have error code: %v (has %v)",
					testCase.epoch, testCase.errorCode, err)
			}
		}
	}
}
