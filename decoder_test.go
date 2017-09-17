package zbase32_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/tv42/zbase32"
)

func TestDecoder(t *testing.T) {
	for _, tc := range byteTests {
		for bs := int64(1); bs < 128; bs += 4 {
			var buf bytes.Buffer
			dec := zbase32.NewDecoder(bytes.NewReader([]byte(tc.encoded)))
			for {
				if _, err := io.CopyN(&buf, dec, bs); io.EOF == err {
					break
				} else if nil != err {
					t.Errorf("Failed to decode: %v", err)
				}
			}
			if g, e := buf.String(), string(tc.decoded); g != e {
				t.Errorf("Decode %x wrong result: %q != %q", tc.decoded, g, e)
			}
		}
	}
}
