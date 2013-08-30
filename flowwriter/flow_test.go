package flowwriter

import (
	"bytes"
	"testing"
)

type Test struct {
	n int
	d string
	w string
}

var ts = []Test{
	{9, `
a b c d e f
  - a b c d
    e f
a b c d e f
`,
		`
a b c d e
f
  - a b c
    d e f
a b c d e
f
`},

	{9, `
a b c d e f
> - a b c d
>   e f
a b c d e f
`,
		`
a b c d e
f
> - a b c
>   d e f
a b c d e
f
`},

	{9, `
a b c d e f
> - a b c d
a b c d e f
`,
		`
a b c d e
f
> - a b c
>   d
a b c d e
f
`},
}

func TestIndent(t *testing.T) {
	for _, ts := range ts {
		b := new(bytes.Buffer)
		w := NewWriter(b, ts.n, "-", "> ", "")
		if _, err := w.Write([]byte(ts.d)); err != nil {
			t.Error(err)
		}
		if err := w.Flush(); err != nil {
			t.Error(err)
		}
		if g := b.String(); g != ts.w {
			t.Logf("%q != %q", g, ts.w)
			t.Log("\n" + g)
			t.Log("\n" + ts.w)
			t.Fail()
		}
	}
}
