package flowwriter

import (
	"bytes"
	"github.com/kr/text/indentwriter"
	"github.com/kr/text/wrapwriter"
	"io"
	"strings"
	"unicode/utf8"
)

type Writer struct {
	width  int
	hang   string
	indent string
	chars   string // hang + indent
	repl rune   // last rune in indent
	noWrap string

	sec  [][]byte
	buf  []byte
	pre0 []byte
	pre  []byte

	w io.Writer
}

func NewWriter(w io.Writer, width int) *Writer {
	nw := new(Writer)
	nw.w = w
	nw.width = width
	nw.hang = "-*"
	nw.indent = "> "
	nw.initChars()
	return nw
}

func (w *Writer) SetHang(chars string) {
	w.hang = chars
	w.initChars()
}

func (w *Writer) SetIndent(chars string) {
	w.indent = chars
	w.initChars()
}

func (w *Writer) initChars() {
	w.chars = w.hang + w.indent
	w.repl, _ = utf8.DecodeLastRuneInString(w.indent)
}

func (w *Writer) SetNoWrap(prefix string) {
	w.noWrap = prefix
}

func (w *Writer) Write(p []byte) (n int, err error) {
	for i, c := range p {
		w.buf = append(w.buf, c)
		if c == '\n' {
			body := bytes.TrimLeft(w.buf, w.chars)
			pre := w.buf[:len(w.buf)-len(body)]
			if !bytes.Equal(pre, w.pre) {
				if err := w.flow(); err != nil {
					return i - len(w.buf), err
				}
				w.pre0 = pre
				w.pre = w.tr(pre)
			}
			w.sec = append(w.sec, body)
			w.buf = nil
		}
	}
	return len(p), nil
}

func (w *Writer) Flush() error {
	// TODO(kr): handle incomplete last line
	return w.flow()
}

func (w *Writer) tr(s []byte) []byte {
	f := func(r rune) rune {
		if strings.ContainsRune(w.hang, r) {
			return w.repl
		}
		return r
	}
	return bytes.Map(f, s)
}

func (w *Writer) flow() (err error) {
	sec := w.sec
	w.sec = nil
	n := w.width - len(w.pre0)
	dst := indentwriter.NewWriter(w.w, [][]byte{w.pre0, w.pre})
	if w.noWrap == "" || string(w.pre) != w.noWrap {
		dst = wrapwriter.NewWriter(dst, n)
	}
	for _, s := range sec {
		if _, err = dst.Write(s); err != nil {
			return err
		}
	}
	if f, ok := dst.(flusher); ok {
		f.Flush()
	}
	return nil
}

func Wrapper(w io.Writer) io.Writer {
	return w
}

type flusher interface {
	Flush() error
}
