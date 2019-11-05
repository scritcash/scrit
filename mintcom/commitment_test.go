package mintcom

import (
	"bytes"
	"io"
	"testing"
)

// TestNewCommitment verifies NewCommitment.
func TestNewCommitment(t *testing.T) {
	// TestData
	tMintID := uint64(1)
	tInput := []byte("Test Input")
	tOutput := []byte("Test Output")
	tProof := []byte("Test Proof")
	tPublicKey := &[PublicKeySize]byte{0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	tPrivateKey := &[PrivateKeySize]byte{0xbd, 0x3c, 0xca, 0xda, 0x6e, 0x56, 0xdb, 0xa7, 0x56, 0x63, 0x81, 0x6c, 0x81, 0x5d, 0x6b, 0x54, 0x2e, 0xb5, 0x0e, 0x80, 0x3b, 0x21, 0x9e, 0x10, 0xbc, 0xdf, 0x9f, 0xe6, 0x66, 0x49, 0x10, 0x13, 0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}

	// Basic operation.
	to, err := NewCommitment(tMintID, tInput, tOutput, tProof, tPublicKey, tPrivateKey)
	if err != nil {
		t.Fatalf("NewCommitment returned unexpected error: %s", err)
	}
	if to.MintID != tMintID {
		t.Error("MintID unset")
	}
	if to.CreateTime <= 1 {
		t.Error("CreateTime unset")
	}
	if to.HO != Hash(tOutput[:]) {
		t.Error("HO unset")
	}
	if to.HP != Hash(tProof[:]) {
		t.Error("HP unset")
	}
	if *to.hi != Hash(tInput[:]) {
		t.Error("hi unset")
	}
	if to.HHI != Hash(to.hi[:]) {
		t.Error("HHI unset")
	}
	if to.HK != HMAC(tPublicKey[:], to.hi[:]) {
		t.Error("HK unset")
	}
	if to.marshalled == nil || len(to.marshalled) != commitmentSize {
		t.Error("marshalled unset")
	}
	if to.VerifySignature(tPublicKey) == false {
		t.Error("unsigned")
	}

	// Random source failure.
	// Prepare readers with deliberate errors.
	origRand := Rand
	defer func() { Rand = origRand }()

	// Reader returns error.
	Rand = ErrorReader{}
	if _, err := NewCommitment(tMintID, tInput, tOutput, tProof, tPublicKey, tPrivateKey); err != io.ErrUnexpectedEOF {
		t.Error("Random source failure not detected")
	}
}

// TestMarshal verifies Marshal/Unmarshal match.
func TestMarshal(t *testing.T) {
	var u3 *Commitment
	// TestData
	tMintID := uint64(1)
	tInput := []byte("Test Input")
	tOutput := []byte("Test Output")
	tProof := []byte("Test Proof")
	tPublicKey := &[PublicKeySize]byte{0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	tPrivateKey := &[PrivateKeySize]byte{0xbd, 0x3c, 0xca, 0xda, 0x6e, 0x56, 0xdb, 0xa7, 0x56, 0x63, 0x81, 0x6c, 0x81, 0x5d, 0x6b, 0x54, 0x2e, 0xb5, 0x0e, 0x80, 0x3b, 0x21, 0x9e, 0x10, 0xbc, 0xdf, 0x9f, 0xe6, 0x66, 0x49, 0x10, 0x13, 0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	td, _ := NewCommitment(tMintID, tInput, tOutput, tProof, tPublicKey, tPrivateKey)

	// Marshal, clear, marshal
	m1 := td.Marshal()
	td.ClearMarshalCache()
	if td.marshalled != nil {
		t.Error("ClearMarshalCache without operation")
	}
	if td.hi == nil {
		t.Error("ClearMarshalCache may not remove .hi")
	}
	m2 := td.Marshal()
	if !bytes.Equal(m1, m2) {
		t.Error("Marshalling is not deterministic")
	}

	// Unmarshal, normal operation on good input.
	u1 := new(Commitment)
	u2 := u1.Unmarshal(m2)
	if u1 != u2 {
		t.Error("Pointer replaced")
	}
	if u2 == nil {
		t.Error("Unmarshal failed")
	}
	u4 := u3.Unmarshal(m2)
	if u4 == nil {
		t.Error("No struct allocated")
	}
	if td.MintID != u2.MintID {
		t.Error("MintID wrong")
	}
	if td.CreateTime != u2.CreateTime {
		t.Error("CreateTime wrong")
	}
	if td.Random != u2.Random {
		t.Error("Random wrong")
	}
	if td.HHI != u2.HHI {
		t.Error("HHI wrong")
	}
	if td.HO != u2.HO {
		t.Error("HO wrong")
	}
	if td.HP != u2.HP {
		t.Error("HP wrong")
	}
	if td.HK != u2.HK {
		t.Error("HK wrong")
	}
	if td.Signature != u2.Signature {
		t.Error("Signature wrong")
	}
	if len(u2.marshalled) != commitmentSize {
		t.Error("marshalled size wrong")
	}
	if !bytes.Equal(m1, u2.marshalled) {
		t.Error("marshalled wrong")
	}
	if u2.hi != nil {
		t.Error("May not set hi")
	}
	u2 = u2.Unmarshal(m2)
	if u2.hi != nil {
		t.Error("Did not reset hi")
	}

	// Test bad input. Wrong package type.
	m2[0] = 0x00
	if u2.Unmarshal(m2) != nil {
		t.Error("Did not detect wrong package type")
	}
	// Short package. Must return 0 and not panic.
	m2 = m2[0 : len(m2)-2]
	if u2.Unmarshal(m2) != nil {
		t.Error("Did not detect short input")
	}
}

// TestSign verifies sign and verifysignature.
func TestSign(t *testing.T) {
	// TestData
	tMintID := uint64(1)
	tInput := []byte("Test Input")
	tOutput := []byte("Test Output")
	tProof := []byte("Test Proof")
	tPublicKey := &[PublicKeySize]byte{0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	tPrivateKey := &[PrivateKeySize]byte{0xbd, 0x3c, 0xca, 0xda, 0x6e, 0x56, 0xdb, 0xa7, 0x56, 0x63, 0x81, 0x6c, 0x81, 0x5d, 0x6b, 0x54, 0x2e, 0xb5, 0x0e, 0x80, 0x3b, 0x21, 0x9e, 0x10, 0xbc, 0xdf, 0x9f, 0xe6, 0x66, 0x49, 0x10, 0x13, 0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	td, _ := NewCommitment(tMintID, tInput, tOutput, tProof, tPublicKey, tPrivateKey)
	td.Signature = [SignatureSize]byte{}
	td.marshalled = nil

	// Normal operation
	td.sign(tPrivateKey)
	td.marshalled = nil
	if td.VerifySignature(tPublicKey) == false {
		t.Error("Verification failed")
	}

	// Modify signature
	td.marshalled = nil
	td.Signature[SignatureSize/2] ^= td.Signature[SignatureSize/2]
	if td.VerifySignature(tPublicKey) == true {
		t.Error("Verification succeeded on wrong signature")
	}
}

//TestVerify verifies Verify() and VerifyCallback()
func TestVerify(t *testing.T) {
	// TestData
	tMintID := uint64(1)
	tInput := []byte("Test Input")
	tInputError := []byte("Test Error Input")
	tHIError := Hash(tInputError)
	tHHIError := Hash(tHIError[:])
	tOutput := []byte("Test Output")
	tProof := []byte("Test Proof")
	tPublicKey := &[PublicKeySize]byte{0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	tPublicKeyFalse := &[PublicKeySize]byte{0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0x14, 0xf3, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	tPrivateKey := &[PrivateKeySize]byte{0xbd, 0x3c, 0xca, 0xda, 0x6e, 0x56, 0xdb, 0xa7, 0x56, 0x63, 0x81, 0x6c, 0x81, 0x5d, 0x6b, 0x54, 0x2e, 0xb5, 0x0e, 0x80, 0x3b, 0x21, 0x9e, 0x10, 0xbc, 0xdf, 0x9f, 0xe6, 0x66, 0x49, 0x10, 0x13, 0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	td, _ := NewCommitment(tMintID, tInput, tOutput, tProof, tPublicKey, tPrivateKey)

	// Normal operation
	if hiok, ok := td.Verify(&td.HHI, td.hi, tPublicKey); hiok != true || ok != true {
		t.Error("Verify failed")
	}
	// Callback, correct key.
	td.marshalled = nil
	cb := func(mintID uint64) *[PublicKeySize]byte {
		if mintID != tMintID {
			t.Fatal("Wrong mintID requested")
		}
		return tPublicKey
	}
	if hiok, ok := td.VerifyLookup(&td.HHI, td.hi, cb); hiok != true || ok != true {
		t.Error("VerifyLookup failed")
	}
	// Callback, wrong key.
	cb = func(mintID uint64) *[PublicKeySize]byte {
		if mintID != tMintID {
			t.Fatal("Wrong mintID requested")
		}
		return tPublicKeyFalse
	}
	if hiok, ok := td.VerifyLookup(&td.HHI, td.hi, cb); hiok != false || ok != false {
		t.Error("VerifyLookup succeeded with wrong key")
	}
	// Callback, no key.
	cb = func(mintID uint64) *[PublicKeySize]byte {
		if mintID != tMintID {
			t.Fatal("Wrong mintID requested")
		}
		return nil
	}
	if hiok, ok := td.VerifyLookup(&td.HHI, td.hi, cb); hiok != false || ok != false {
		t.Error("VerifyLookup succeeded with no key")
	}

	// only input
	if hiok, ok := td.Verify(nil, td.hi, tPublicKey); hiok != true || ok != true {
		t.Error("Verify failed only hi")
	}
	// only hhi
	if hiok, ok := td.Verify(&td.HHI, nil, tPublicKey); hiok != false || ok != true {
		t.Error("Verify failed only hhi")
	}

	// Wrong HHI
	if hiok, ok := td.Verify(&tHHIError, nil, tPublicKey); hiok != false || ok != false {
		t.Error("Verify failed wrong hhi")
	}
	if hiok, ok := td.Verify(&tHHIError, td.hi, tPublicKey); hiok != false || ok != false {
		t.Error("Verify failed wrong hhi with hi")
	}
	if hiok, ok := td.Verify(nil, &tHIError, tPublicKey); hiok != false || ok != false {
		t.Error("Verify failed wrong hi no hhi")
	}
	if hiok, ok := td.Verify(&tHHIError, &tHIError, tPublicKey); hiok != false || ok != false {
		t.Error("Verify failed wrong hi and hhi")
	}

	// Wrong HI
	if hiok, ok := td.Verify(&td.HHI, &tHIError, tPublicKey); hiok != false || ok != true {
		t.Errorf("Verify failed wrong hi. %t %t", hiok, ok)
	}

	// Modify marshal value. ToDo: Test order
	td.Marshal()
	td.marshalled[0] = 0x00
	if hiok, ok := td.Verify(nil, td.hi, tPublicKey); hiok != false || ok != false {
		t.Error("PkgType error undetected")
	}
}

// TestMatches verifies the Match function
func TestMatches(t *testing.T) {
	// TestData
	tMintID := uint64(1)
	tProof := []byte("Test Proof")
	tPublicKey := &[PublicKeySize]byte{0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}
	tPrivateKey := &[PrivateKeySize]byte{0xbd, 0x3c, 0xca, 0xda, 0x6e, 0x56, 0xdb, 0xa7, 0x56, 0x63, 0x81, 0x6c, 0x81, 0x5d, 0x6b, 0x54, 0x2e, 0xb5, 0x0e, 0x80, 0x3b, 0x21, 0x9e, 0x10, 0xbc, 0xdf, 0x9f, 0xe6, 0x66, 0x49, 0x10, 0x13, 0xd4, 0x4b, 0xda, 0x03, 0x6c, 0x56, 0x72, 0xf0, 0xcc, 0x2a, 0x0e, 0xc5, 0x4b, 0xc8, 0x3f, 0x1a, 0xd0, 0xf3, 0x14, 0x94, 0xc0, 0xdc, 0xec, 0xea, 0xa1, 0xaf, 0x8e, 0xbb, 0xdb, 0x3e, 0x42, 0x95}

	t1Input := []byte("Test Input")
	t1Output := []byte("Test Output")
	td1, _ := NewCommitment(tMintID, t1Input, t1Output, tProof, tPublicKey, tPrivateKey)

	// Same
	t2Input := []byte("Test Input")
	t2Output := []byte("Test Output")
	td2, _ := NewCommitment(tMintID, t2Input, t2Output, tProof, tPublicKey, tPrivateKey)

	if same, ok := td1.Matches(td2); same != true || ok != true {
		t.Error("Match not recognized")
	}

	// HHI different
	t3Input := []byte("Test Input BAD")
	t3Output := []byte("Test Output")
	td3, _ := NewCommitment(tMintID, t3Input, t3Output, tProof, tPublicKey, tPrivateKey)

	if same, ok := td1.Matches(td3); same != false || ok != false {
		t.Error("Unequal input not recognized")
	}

	// Output different
	t4Input := []byte("Test Input")
	t4Output := []byte("Test Output BAD")
	td4, _ := NewCommitment(tMintID, t4Input, t4Output, tProof, tPublicKey, tPrivateKey)

	if same, ok := td1.Matches(td4); same != true || ok != false {
		t.Error("Unequal output not recognized")
	}

	// HHI and output different
	t5Input := []byte("Test Input BAD")
	t5Output := []byte("Test Output BAD")
	td5, _ := NewCommitment(tMintID, t5Input, t5Output, tProof, tPublicKey, tPrivateKey)
	if same, ok := td1.Matches(td5); same != false || ok != false {
		t.Error("Unequal input and output not recognized")
	}

}
