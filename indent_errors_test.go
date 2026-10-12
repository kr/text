package text

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type interruptedWriter struct {
	bytes.Buffer
	call, failCall, accepted int
	err                      error
}

func (w *interruptedWriter) Write(p []byte) (int, error) {
	w.call++
	if w.call == w.failCall {
		w.Buffer.Write(p[:w.accepted])
		return w.accepted, w.err
	}
	return w.Buffer.Write(p)
}

func TestIndentWriterInterruptedWrite(t *testing.T) {
	failure := errors.New("write interrupted")
	for _, tc := range []struct {
		name                  string
		call, accepted, wantN int
		err, wantErr          error
	}{
		{"data accepted with error", 2, 1, 1, failure, failure},
		{"newline accepted with error", 3, 1, 2, failure, failure},
		{"partial prefix with error", 1, 1, 0, failure, failure},
		{"short prefix without error", 1, 1, 0, nil, io.ErrShortWrite},
		{"short data without error", 2, 0, 0, nil, io.ErrShortWrite},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output := &interruptedWriter{failCall: tc.call, accepted: tc.accepted, err: tc.err}
			writer := NewIndentWriter(output, []byte("> "), []byte("- "))
			input := []byte("a\nb")
			n, err := writer.Write(input)
			if n != tc.wantN || err != tc.wantErr {
				t.Fatalf("Write = (%d, %v), want (%d, %v)", n, err, tc.wantN, tc.wantErr)
			}
			if _, err := writer.Write(input[n:]); err != nil {
				t.Fatal(err)
			}
			if got := output.String(); got != "> a\n- b" {
				t.Errorf("resumed output = %q", got)
			}
		})
	}
}
