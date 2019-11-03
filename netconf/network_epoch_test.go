package netconf

import (
	"testing"
)

func TestEpochValidate(t *testing.T) {
	testCases := []struct {
		epoch     NetworkEpoch
		errorCode error
	}{
		{
			NetworkEpoch{
				M: 0,
				N: 1,
			},
			ErrZeroM,
		},
		{
			NetworkEpoch{
				M: 1,
				N: 0,
			},
			ErrZeroN,
		},
		{
			NetworkEpoch{
				M: 2,
				N: 1,
			},
			ErrMGreaterN,
		},
		{
			NetworkEpoch{
				M: 5,
				N: 10,
			},
			ErrQuorumTooSmall,
		},
		{
			NetworkEpoch{
				M: 6,
				N: 10,
			},
			nil,
		},
		{
			NetworkEpoch{
				M: 5,
				N: 11,
			},
			ErrQuorumTooSmall,
		},
		{
			NetworkEpoch{
				M: 6,
				N: 11,
			},
			nil,
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
