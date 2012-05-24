package text

import (
	"testing"
)

type T struct {
	inp, exp, pre string
}

var tests = []T{
	{
		"The quick brown fox\njumps over the lazy\ndog.\nBut not quickly.\n",
		"xxxThe quick brown fox\nxxxjumps over the lazy\nxxxdog.\nxxxBut not quickly.\n",
		"xxx",
	},
	{
		"The quick brown fox\njumps over the lazy\ndog.\n\nBut not quickly.",
		"xxxThe quick brown fox\nxxxjumps over the lazy\nxxxdog.\n\nxxxBut not quickly.",
		"xxx",
	},
}

func TestIndent(t *testing.T) {
	for _, test := range tests {
		got := Indent(test.inp, test.pre)
		if got != test.exp {
			t.Errorf("mismatch %q != %q", got, test.exp)
		}
	}
}
