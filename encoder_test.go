package zbase32_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/tv42/zbase32"
)

func TestEncoder(t *testing.T) {
	for _, tc := range byteTests {
		for bs := int64(1); bs < 128; bs += 4 {
			in := bytes.NewReader(tc.decoded)
			buf := new(bytes.Buffer)
			enc := zbase32.NewEncoder(buf)
			for {
				if _, err := io.CopyN(enc, in, bs); io.EOF == err {
					break
				} else if nil != err {
					t.Errorf("Failed to encode: %v", err)
				}
			}
			if err := enc.Close(); nil != err {
				t.Errorf("Failed to close encoder: %v", err)
			}

			if g, e := buf.String(), tc.encoded; g != e {
				t.Errorf("Encode %x wrong result: %q != %q", tc.decoded, g, e)
				continue
			}
		}
	}

}
