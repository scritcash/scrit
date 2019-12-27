/*
Package binencode provides functions to encode values into length-encoded
slices. It can be used to work on secure memory. Only int16, int32, int64 and
[]byte are supported. Inserting and skipping zero bytes is supported.

Allocation-Free Serialization

This library supports efficient serialization/deserialization of most relevant
data types into byte slices. Message type and data type encoding are
supported, including variable length encodings. No support for stacked
structs exists or will be included.

It uses only pre-allocated memory so that manual allocation and thus memory
protection schemes are accessible.

Two methods for variable allocation are availble: With reflection and without
reflection, based only on language type annotations. This simplifies
development (reflect) while making it very easy to switch to efficient
runtime.

  type data struct {
    a int16
    b int32
    c int64
    d []byte
  }

  // ...

  v := &data{
    a: 1,
    b: 29381,
    c: 5098123,
    d: []byte("any string"),
     }

  // Define fields to include and their order in the output message.
  encodingScheme := []interface{}{v.a, v.b, v.c, v.d}

  // Calculate size of encoded message for manual allocation.
  size := EncodeSize(encodingScheme)

  // Serialize data into message using preallocated buffer and the defined
  // scheme.
  encodedData, _ := Encode(buf, encodingScheme...)

  // Decode message into the pre-allocated variable decribed by encodingScheme.
  _,_ := Decode(encodedData, encodingScheme)
*/
package binencode
