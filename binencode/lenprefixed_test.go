package binencode

import (
	"bytes"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestInt16(t *testing.T) {
	var ta, tb int16
	d := make([]byte, 0, 10000)
	s := 0
	b, l, err := EncodeInt16(15, d)
	if err != nil {
		t.Errorf("EncodeInt16: %s", err)
	}
	s += l
	b, l, err = EncodeInt16(3, b)
	if err != nil {
		t.Errorf("EncodeInt16: %s", err)
	}
	s += l
	b, l, err = DecodeInt16(d[0:s], &ta)
	if err != nil {
		t.Errorf("DecodeInt16: %s", err)
	}
	b, l, err = DecodeInt16(b, &tb)
	if err != nil {
		t.Errorf("DecodeInt16: %s", err)
	}
	if ta != 15 || tb != 3 {
		t.Error("Decoded values don't match")
	}
}

func TestInt32(t *testing.T) {
	var ta, tb int32
	d := make([]byte, 0, 10000)
	s := 0
	b, l, err := EncodeInt32(15, d)
	if err != nil {
		t.Errorf("EncodeInt32: %s", err)
	}
	s += l
	b, l, err = EncodeInt32(3, b)
	if err != nil {
		t.Errorf("EncodeInt32: %s", err)
	}
	s += l
	b, l, err = DecodeInt32(d[0:s], &ta)
	if err != nil {
		t.Errorf("DecodeInt32: %s", err)
	}
	b, l, err = DecodeInt32(b, &tb)
	if err != nil {
		t.Errorf("DecodeInt32: %s", err)
	}
	if ta != 15 || tb != 3 {
		t.Error("Decoded values don't match")
	}
}

func TestInt64(t *testing.T) {
	var ta, tb int64
	d := make([]byte, 0, 10000)
	s := 0
	b, l, err := EncodeInt64(15, d)
	if err != nil {
		t.Errorf("EncodeInt64: %s", err)
	}
	s += l
	b, l, err = EncodeInt64(3, b)
	if err != nil {
		t.Errorf("EncodeInt64: %s", err)
	}
	s += l
	b, l, err = DecodeInt64(d[0:s], &ta)
	if err != nil {
		t.Errorf("DecodeInt64: %s", err)
	}
	b, l, err = DecodeInt64(b, &tb)
	if err != nil {
		t.Errorf("DecodeInt64: %s", err)
	}
	if ta != 15 || tb != 3 {
		t.Error("Decoded values don't match")
	}
}

func TestBytes(t *testing.T) {
	td := []byte("test data 1234")
	s := EncodeBytesSize(td)
	if s != len(td)+5 {
		t.Error("EncodeBytesSize bad value")
	}
	tmp := make([]byte, 0, EncodeBytesSize(td)+2)
	_, n, err := EncodeBytes(td, tmp)
	if err != nil {
		t.Errorf("EncodeBytes: %s", err)
	}
	if n != s {
		t.Error("Encoding wrong size")
	}
	out := tmp[0:cap(tmp)]
	ds, ok := DecodeBytesSize(out)
	if !ok {
		t.Error("DecodeBytesSize not ok")
	}
	if ds != len(td) {
		t.Errorf("DecodeBytesSize wrong size: %d != %d", len(td), ds)
	}
	var td2 []byte
	if _, _, err := DecodeBytes(out, &td2); err != nil {
		t.Errorf("DecodeBytes: %s", err)
	}
	if !bytes.Equal(td, td2) {
		t.Error("Decode failure, nil")
	}
	td3 := make([]byte, 0, ds+2)
	if _, _, err := DecodeBytes(out, &td3); err != nil {
		t.Errorf("DecodeBytes: %s", err)
	}
	if !bytes.Equal(td, td3) {
		spew.Dump(td3)
		t.Error("Decode failure, cap")
	}
	td4 := make([]byte, ds)
	if _, _, err := DecodeBytes(out, &td4); err != nil {
		t.Errorf("DecodeBytes: %s", err)
	}
	if !bytes.Equal(td, td4) {
		t.Error("Decode failure, len")
	}
	td5 := make([]byte, 3)
	if _, _, err := DecodeBytes(out, &td5); err == nil {
		t.Error("DecodeBytes Len short")
	}
	td5 = make([]byte, 100)
	if _, _, err := DecodeBytes(out, &td5); err == nil {
		t.Error("DecodeBytes Len excess")
	}
	td5 = make([]byte, 0, 3)
	if _, _, err := DecodeBytes(out, &td5); err == nil {
		t.Error("DecodeBytes Cap short")
	}
	var td6 [14]byte
	if _, _, err := DecodeBytes(out, SlicePointer(td6[:])); err != nil {
		t.Errorf("DecodeBytes: %s", err)
	}
	if !bytes.Equal(td, td6[:]) {
		t.Error("Decode failure, array")
	}
}

func TestBytesZero(t *testing.T) {
	td := make([]byte, 0)
	s := EncodeBytesSize(td)
	if s != len(td)+5 {
		t.Error("EncodeBytesSize bad value")
	}
	tmp := make([]byte, 0, EncodeBytesSize(td)+2)
	_, _, err := EncodeBytes(td, tmp)
	if err != nil {
		t.Errorf("EncodeBytes: %s", err)
	}
	out := tmp[0:cap(tmp)]
	var td2 []byte
	if _, _, err := DecodeBytes(out, &td2); err != nil {
		t.Errorf("DecodeBytes: %s", err)
	}
}

func TestBytesNil(t *testing.T) {
	var td []byte
	s := EncodeBytesSize(td)
	if s != len(td)+5 {
		t.Error("EncodeBytesSize bad value")
	}
	tmp := make([]byte, 0, EncodeBytesSize(td)+2)
	_, _, err := EncodeBytes(td, tmp)
	if err != nil {
		t.Errorf("EncodeBytes: %s", err)
	}
	out := tmp[0:cap(tmp)]
	var td2 []byte
	if _, _, err := DecodeBytes(out, &td2); err != nil {
		t.Errorf("DecodeBytes: %s", err)
	}
}

func TestSkip(t *testing.T) {
	td := []byte("1234567890")
	_, _, err := EncodeSkip(5, td[:0])
	if err != nil {
		t.Errorf("EncodeSkip: %s", err)
	}
	DecodeSkip(td[0:cap(td)], 5)
}

type testStruct struct {
	a int16
	b int32
	c int64
	d []byte
}

type testStruct2 struct {
	A int16
	B int32
	C int64
	D []byte
	c int32
}

func TestReflect(t *testing.T) {
	td := &testStruct2{
		A: 3,
		B: 15,
		C: 1239123,
		D: []byte("test value"),
	}
	td2 := new(testStruct2)
	desc := DescribeStruct(td)
	desc2 := DescribeStruct(td2)
	out, err := Encode(nil, desc...)
	if err != nil {
		t.Fatalf("Encode: %s", err)
	}
	_, err = Decode(out, desc2...)
	if err != nil {
		t.Errorf("Decode: %s", err)
	}
	if td.A != td2.A || td.B != td2.B || td.C != td2.C || !bytes.Equal(td.D, td2.D) {
		t.Error("Decode corrupt")
	}
}

func TestConvenience(t *testing.T) {
	td := &testStruct{
		a: 1,
		b: 2,
		c: 3,
		d: []byte("test"),
	}
	td2 := new(testStruct)
	val1 := []interface{}{td.a, 3, td.b, td.c, td.d}
	val2 := []interface{}{&td.a, 3, &td.b, &td.c, &td.d}
	val3 := []interface{}{&td2.a, 3, &td2.b, &td2.c, &td2.d}
	_, err := EncodeSize(val1...)
	if err != nil {
		t.Fatalf("EncodeSize: %s", err)
	}
	_, err = EncodeSize(val2...)
	if err != nil {
		t.Fatalf("EncodeSize Pointers: %s", err)
	}
	out, err := Encode(nil, val1...)
	if err != nil {
		t.Fatalf("Encode: %s", err)
	}
	out2, err := Encode(nil, val2...)
	if err != nil {
		t.Fatalf("Encode Pointers: %s", err)
	}
	if !bytes.Equal(out, out2) {
		t.Error("Encoded values should match")
	}
	rem, err := Decode(out, val3...)
	if err != nil {
		t.Errorf("Decode: %s", err)
	}
	if len(rem) != 0 {
		t.Error("There should be no remainder")
	}
	if td.a != td2.a || td.b != td2.b || td.c != td2.c || !bytes.Equal(td.d, td2.d) {
		t.Error("Decode corrupt")
	}
}
