package wrapwriter

import (
	"bytes"
	"testing"
)

const text = "The quick brown fox jumps over the lazy dog."

type Test struct {
	n   int
	in  string
	exp string
}

var ts = []Test{
	{24, text, "The quick brown fox\njumps over the lazy dog.\n"},
	{5, text, "The\nquick\nbrown\nfox\njumps\nover\nthe\nlazy\ndog.\n"},
	{500, text, "The quick brown fox jumps over the lazy dog.\n"},
	{9, "\na b c d e f\n", "\na b c d e\nf\n"},
	{9, "\na b c d e f\n\na b c d e f\n", "\na b c d e\nf\n\na b c d e\nf\n"},
}

func TestWrap(t *testing.T) {
	for _, ts := range ts {
		b := new(bytes.Buffer)
		w := NewWriter(b, ts.n)
		if _, err := w.Write([]byte(ts.in)); err != nil {
			t.Error(err)
		}
		if err := w.Flush(); err != nil {
			t.Error(err)
		}
		got := b.String()
		if got != ts.exp {
			t.Errorf("%q != %q", got, ts.exp)
			t.Log(got)
			t.Log(ts.exp)
		}
	}
}
