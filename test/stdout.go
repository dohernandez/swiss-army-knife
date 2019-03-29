package test

import (
	"bytes"
	"io"
	"os"
)

// CaptureStdOut returns whatever is print to Stdout in the scope of the function pass as a parameter.
func CaptureStdOut(f func()) string {
	old := os.Stdout // keep backup of the real stdout.
	// nolint:errcheck
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely.
	go func() {
		var buf bytes.Buffer
		// nolint:errcheck
		// #nosec G104
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state.
	// nolint:errcheck
	// #nosec G104
	w.Close()
	os.Stdout = old // restoring the real stdout.
	out := <-outC

	return out
}
