package zbase32_test

import (
	"fmt"
	"os"

	"github.com/tv42/zbase32"
)

func Example() {
	s := zbase32.EncodeToString([]byte{240, 191, 199})
	fmt.Println(s)
	// Output:
	// 6n9hq
}

func ExampleNewEncoder() {
	input := []byte("foo\x00bar")
	encoder := zbase32.NewEncoder(os.Stdout)
	encoder.Write(input)
	// Must close the encoder when finished to flush any partial blocks.
	encoder.Close()
	// Output:
	// c3zs6ydncf3y
}
