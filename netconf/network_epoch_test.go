package netconf

import (
	"testing"
	"time"
)

var (
	t1 time.Time
	t2 time.Time
	t3 time.Time
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
}

func TestEpochValidate(t *testing.T) {
	testCases := []struct {
		epoch     NetworkEpoch
		errorCode error
	}{
		{
			NetworkEpoch{
				M:           0,
				N:           1,
				SignStart:   t1,
				SignEnd:     t2,
				ValidateEnd: t3,
			},
			ErrZeroM,
		},
		{
			NetworkEpoch{
				M:           1,
				N:           0,
				SignStart:   t1,
				SignEnd:     t2,
				ValidateEnd: t3,
			},
			ErrZeroN,
		},
		{
			NetworkEpoch{
				M:           2,
				N:           1,
				SignStart:   t1,
				SignEnd:     t2,
				ValidateEnd: t3,
			},
			ErrMGreaterN,
		},
		{
			NetworkEpoch{
				M:           5,
				N:           10,
				SignStart:   t1,
				SignEnd:     t2,
				ValidateEnd: t3,
			},
			ErrQuorumTooSmall,
		},
		{
			NetworkEpoch{
				M:           6,
				N:           10,
				SignStart:   t1,
				SignEnd:     t2,
				ValidateEnd: t3,
			},
			nil,
		},
		{
			NetworkEpoch{
				M:           5,
				N:           11,
				SignStart:   t1,
				SignEnd:     t2,
				ValidateEnd: t3,
			},
			ErrQuorumTooSmall,
		},
		{
			NetworkEpoch{
				M:           6,
				N:           11,
				SignStart:   t1,
				SignEnd:     t2,
				ValidateEnd: t3,
			},
			nil,
		},
		{
			NetworkEpoch{
				M:         8,
				N:         10,
				SignStart: t1,
				SignEnd:   t1,
			},
			ErrSignEpochStartNotBeforeSignEnd,
		},
		{
			NetworkEpoch{
				M:           8,
				N:           10,
				SignStart:   t1,
				SignEnd:     t2,
				ValidateEnd: t2,
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
