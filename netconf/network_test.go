package netconf

import (
	"path/filepath"
	"testing"
)

func TestLoadNetwork(t *testing.T) {
	net, err := LoadNetwork(filepath.Join("testdata", DefNetConfFile))
	if err != nil {
		t.Fatal(err)
	}
	if err := net.Validate(); err != nil {
		t.Fatal(err)
	}
	_ = DBCTypeMapToSortedArray(net.DBCTypes())
}

func TestNetworkValidate(t *testing.T) {
	testCases := []struct {
		net       Network
		errorCode error
	}{
		{
			Network{
				[]NetworkEpoch{
					{
						QuorumM:        8,
						NumberOfMintsN: 10,
						SignStart:      t1,
						SignEnd:        t2,
						ValidateEnd:    t3,
					},
					{
						QuorumM:        8,
						NumberOfMintsN: 10,
						SignStart:      t1,
						SignEnd:        t2,
						ValidateEnd:    t3,
					},
				},
			},
			ErrSignEpochWrongBoundaries,
		},
		{
			Network{
				[]NetworkEpoch{
					{
						QuorumM:        8,
						NumberOfMintsN: 10,
						SignStart:      t1,
						SignEnd:        t2,
						ValidateEnd:    t4,
					},
					{
						QuorumM:        8,
						NumberOfMintsN: 10,
						SignStart:      t2,
						SignEnd:        t3,
						ValidateEnd:    t4,
					},
				},
			},
			ErrValidationLongerThanNextSigning,
		},
	}
	for _, testCase := range testCases {
		err := testCase.net.Validate()
		if err != testCase.errorCode {
			if err != testCase.errorCode {
				t.Fatalf("Validate(%#v) should have error code: %v (has %v)",
					testCase.net, testCase.errorCode, err)
			}
		}
	}
}
