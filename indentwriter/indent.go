package indentwriter

import (
	"io"
)

// Writer indents each line of its input.
type Writer struct {
	w   io.Writer
	bol bool
	pre [][]byte
	sel int
}

// NewWriter makes a new write filter that indents the input lines.
// Each line is prefixed with the corresponding element of pre. If
// there are more lines than elements, the last element of pre is
// repeated for each subsequent line.
func NewWriter(w io.Writer, pre [][]byte) io.Writer {
	return &Writer{
		w:   w,
		pre: pre,
		bol: true,
	}
}

// The only errors returned are from the underlying writer.
func (w *Writer) Write(p []byte) (n int, err error) {
	for _, c := range p {
		if w.bol {
			if _, err = w.w.Write(w.pre[w.sel]); err != nil {
				return n, err
			}
		}
		_, err = w.w.Write([]byte{c})
		if err != nil {
			return n, err
		}
		n++
		w.bol = c == '\n'
		if w.bol && w.sel < len(w.pre)-1 {
			w.sel++
		}
	}
	return n, nil
}
