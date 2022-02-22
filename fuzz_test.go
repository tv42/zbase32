//go:build go1.18
// +build go1.18

package zbase32_test

import (
	"encoding/binary"
	"testing"

	"github.com/tv42/zbase32"
)

func FuzzDecodeString(f *testing.F) {
	f.Add("")
	f.Add("y")
	f.Add("o")
	f.Add("6im54d")
	f.Add("9h")
	f.Add("99o")
	f.Add("99999999")
	f.Add("999999999h")
	f.Add("ab3sr1ix8fhfnuzaeo75fkn3a7xh8udk6jsiiko")
	f.Fuzz(func(t *testing.T, input string) {
		decoded, err := zbase32.DecodeString(input)
		// For fuzzing, we don't care whether it parses right or not.
		// testing.F has no mechanism to return "interestingness", otherwise we'd signal err == nil as interesting.
		_ = decoded
		_ = err
	})
}

func FuzzDecodeBitsString(f *testing.F) {
	buf := make([]byte, 8)
	for num := uint64(0); num < 1_000_000; num += 1000 {
		binary.BigEndian.PutUint64(buf, num)
		for bits := 0; bits < 14; bits++ {
			output := zbase32.EncodeBitsToString(buf, bits)
			f.Add(output, bits)
		}
	}
	f.Fuzz(func(t *testing.T, input string, bits int) {
		decoded, err := zbase32.DecodeBitsString(input, bits)
		// For fuzzing, we don't care whether it parses right or not.
		// testing.F has no mechanism to return "interestingness", otherwise we'd signal err == nil as interesting.
		_ = decoded
		_ = err
	})
}

func FuzzEncodeBits(f *testing.F) {
	f.Add([]byte("foobar"), 17)
	f.Fuzz(func(t *testing.T, input []byte, bits int) {
		if bits > len(input)*8 {
			return
		}
		dst := make([]byte, zbase32.EncodedLen(len(input)))
		err := zbase32.EncodeBits(dst, input, bits)
		_ = err
	})
}
