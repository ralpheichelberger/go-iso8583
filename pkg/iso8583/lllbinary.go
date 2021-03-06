package iso8583

import (
	"errors"
)

// LLLBINARY is a []byte implementation of a field with a LLL indicator before which can be encoded using encode tag,
// it does not contain any special behaviour more than unload all bytes on marshaling and
// reading the specified length on unmarshaling.
type LLLBINARY []byte

// MarshalISO8583 returns a copy of binary content. Encoding and length input are ignored.
func (binary LLLBINARY) MarshalISO8583(length int, enc string) ([]byte, error) {
	binaryCopy := make([]byte, len(binary))
	copy(binaryCopy, binary)

	_, llEncoding := ReadSplitEncodings(enc)
	return LengthMarshal(3, binaryCopy, llEncoding)
}

// UnmarshalISO8583 reads the length indicated amount of bytes from b and load the BINARY field with it.
// Encoding is ignored.
func (binary *LLLBINARY) UnmarshalISO8583(b []byte, length int, enc string) (int, error) {
	if b == nil {
		return 0, errors.New("bytes input is nil")
	}

	lllEncoding, _ := ReadSplitEncodings(enc)

	n, b, err := LengthUnmarshal(3, b, length, lllEncoding)
	if err != nil {
		return 0, err
	}

	*binary = make([]byte, len(b))
	copy(*binary, b)

	return n, nil
}
