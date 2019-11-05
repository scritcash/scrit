package mintcom

// %%TODO%%: Always include proof in commitment.

import (
	"crypto/hmac"
	"encoding/binary"

	"golang.org/x/crypto/ed25519"
)

// RandomSize is the number of random bytes in a commitment.
const RandomSize = HashSize / 2

// SignatureSize is the size of signatures used.
const SignatureSize = ed25519.SignatureSize

// PublicKeySize is the size of public keys for the signature algorithm.
const PublicKeySize = ed25519.PublicKeySize

// PrivateKeySize is the size of private keys for the signature algorithm.
const PrivateKeySize = ed25519.PrivateKeySize

// packageSize: Type(1byte)+MintID+CreateTime+Random+HHI+HO+HP+HK
const packageSize = 1 + 8 + 8 + RandomSize + HashSize + HashSize + HashSize + HashSize

// commitmentSize: packageSize+SignatureSize
const commitmentSize = packageSize + SignatureSize

// Commitment contains an input:output commitment by a mint.
type Commitment struct {
	MintID     uint64              // The public ID of the Mint.
	CreateTime uint64              // The time when the commitment was created.
	Random     [RandomSize]byte    // Random bytes to protect NotFound replies.
	HHI        [HashSize]byte      // Hash(Hash(Input)).
	HO         [HashSize]byte      // Hash(Output).
	HP         [HashSize]byte      // Hash(Proof).
	HK         [HashSize]byte      // HMAC(MintPubkey, Hash(Input))
	Signature  [SignatureSize]byte // Signature over the above.

	hi         *[HashSize]byte // Hash(Input)
	marshalled []byte          // Marshalled commitment.
}

// NewCommitment creates a new commitment from the given parameters.
func NewCommitment(mintID uint64, input, output, proof []byte, publicKey *[PublicKeySize]byte, privateKey *[PrivateKeySize]byte) (*Commitment, error) {
	com := &Commitment{
		MintID:     mintID,
		CreateTime: Now(),
		HO:         Hash(output),
		HP:         Hash(proof),
	}
	if err := RandomBytes(com.Random[:]); err != nil {
		return nil, err
	}
	hi := Hash(input)
	com.hi = &hi
	com.HHI = Hash(com.hi[:])
	com.HK = HMAC(publicKey[:], com.hi[:])
	com.sign(privateKey)
	return com, nil
}

func (com *Commitment) sign(privateKey *[PrivateKeySize]byte) {
	if com.marshalled == nil {
		com.Marshal()
	}
	sig := ed25519.Sign(privateKey[:], com.marshalled[0:packageSize])
	copy(com.Signature[:], sig)
	com.marshalled = com.marshalled[0:commitmentSize]
	copy(com.marshalled[packageSize:commitmentSize], com.Signature[:])
}

// VerifySignature verifies the signature on a commitment.
func (com *Commitment) VerifySignature(publicKey *[PublicKeySize]byte) bool {
	if com.marshalled == nil {
		com.Marshal()
	}
	return ed25519.Verify(publicKey[:], com.marshalled[0:packageSize], com.marshalled[packageSize:commitmentSize])
}

// Verify verifies the commitment. It tests for commitment on input hi if input is not nil. hhi is the Hash(Hash(input)).
// It returns hiok == true if hi was present and verified, false otherwise. ok is returned if all executed tests verify.
func (com *Commitment) Verify(hhi *[HashSize]byte, hi *[HashSize]byte, publicKey *[PublicKeySize]byte) (hiok bool, ok bool) {
	if com.marshalled == nil {
		com.Marshal()
	}
	// Set hhi if hhi is not set but hi is.
	if hhi == nil && hi != nil {
		c := Hash(hi[:])
		hhi = &c
	}
	// Verify package type. Only meaningful if executed on commitment with existing cache.
	if com.marshalled[0] != PkgTypeCommitment {
		return false, false
	}
	// Verify that commitment uses correct hhi.
	if com.HHI != *hhi {
		return false, false
	}
	// Verify signature on commitment.
	if com.VerifySignature(publicKey) == false {
		return false, false
	}
	if hi != nil {
		ht := HMAC(publicKey[:], hi[:])
		if hmac.Equal(ht[:], com.HK[:]) {
			com.hi = hi
			return true, true
		}
		return false, true
	}
	return false, true
}

// VerifyLookup verifies the commitment. It tests for commitment on input hi if input is not nil. hhi is the Hash(Hash(input)).
// It returns hiok == true if hi was present and verified, false otherwise. ok is returned if all executed tests verify.
func (com *Commitment) VerifyLookup(hhi *[HashSize]byte, hi *[HashSize]byte, keyLookup PublicKeyLookup) (hiok bool, ok bool) {
	publicKey := keyLookup(com.MintID)
	if publicKey == nil {
		return false, false
	}
	return com.Verify(hhi, hi, publicKey)
}

// Matches tests if the parameter commitment d matches the receiver commitment com. It does not perform verification of the commitments.
// It returns same==false if the commitments are for different outputs. It returns ok==true if the commitments match.
func (com *Commitment) Matches(d *Commitment) (same, ok bool) {
	if com.HHI != d.HHI {
		return false, false
	}
	if com.HO != d.HO {
		return true, false
	}
	return true, true
}

// ClearMarshalCache clears the internal marshalling cache.
func (com *Commitment) ClearMarshalCache() {
	com.marshalled = nil
}

// Marshal the Commitment into a byte slice. Returns the cached marshalled value if available.
func (com *Commitment) Marshal() []byte {
	if com.marshalled != nil {
		return com.marshalled
	}
	com.marshalled = make([]byte, 1+8+8, commitmentSize) // 1+8+8 is the size of Type||MintID||CreateTime
	com.marshalled[0] = PkgTypeCommitment
	binary.BigEndian.PutUint64(com.marshalled[1:1+8], com.MintID)
	binary.BigEndian.PutUint64(com.marshalled[1+8:1+8+8], com.CreateTime)
	com.marshalled = append(com.marshalled, com.Random[:]...)
	com.marshalled = append(com.marshalled, com.HHI[:]...)
	com.marshalled = append(com.marshalled, com.HO[:]...)
	com.marshalled = append(com.marshalled, com.HP[:]...)
	com.marshalled = append(com.marshalled, com.HK[:]...)
	com.marshalled = append(com.marshalled, com.Signature[:]...)
	return com.marshalled
}

// Unmarshal d into a commitment struct. Writes into receiver and returns the result. Returns nil if unmarshalling is unsuccessful.
// If receiver is nil, a new Commitment struct is allocated.
func (com *Commitment) Unmarshal(d []byte) *Commitment {
	var r *Commitment
	if len(d) < commitmentSize {
		return nil
	}
	if d[0] != PkgTypeCommitment {
		return nil
	}
	r = com
	if com == nil {
		r = new(Commitment)
	}
	r.MintID = binary.BigEndian.Uint64(d[1 : 1+8])
	r.CreateTime = binary.BigEndian.Uint64(d[1+8 : 1+8+8])
	copy(r.Random[:], d[1+8+8:1+8+8+RandomSize])
	copy(r.HHI[:], d[1+8+8+RandomSize:1+8+8+RandomSize+HashSize])
	copy(r.HO[:], d[1+8+8+RandomSize+HashSize:1+8+8+RandomSize+HashSize+HashSize])
	copy(r.HP[:], d[1+8+8+RandomSize+HashSize+HashSize:1+8+8+RandomSize+HashSize+HashSize+HashSize])
	copy(r.HK[:], d[1+8+8+RandomSize+HashSize+HashSize+HashSize:1+8+8+RandomSize+HashSize+HashSize+HashSize+HashSize])
	copy(r.Signature[:], d[1+8+8+RandomSize+HashSize+HashSize+HashSize+HashSize:1+8+8+RandomSize+HashSize+HashSize+HashSize+HashSize+SignatureSize])
	r.marshalled = make([]byte, commitmentSize)
	copy(r.marshalled, d[0:commitmentSize])
	r.hi = nil
	return r
}
