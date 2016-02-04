package text

import (
	"bytes"
	"math"
	"unicode/utf8"
)

var (
	nl = []byte{'\n'}
	sp = []byte{' '}
)

const defaultPenalty = 1e5

// Wrapper splits a list of words into lines with minimal
// "raggedness", attempting to limit lines to MaxWidth units.
//
// Raggedness is the total error over all lines, where error is
// the square of the difference of the length of the line and
// MaxWidth. Too-long lines have an extra penalty added to the
// error -- that is, it is much better for a line to be too short
// than too long.
//
// The zero value is a valid wrapper, ready to use.
type Wrapper struct {
	// MaxWidth is the target for the maximum width of each line
	// If 0, it uses 65 * Width(' ').
	MaxWidth int

	// Penalty is the extra penalty applied to the error
	// for lines that would be too long.
	// If 0, it uses 1e5 * Width(' ').
	Penalty int

	// Width computes the width (in arbitrary units)
	// of the text in b.
	// If nil, it uses utf8.RuneCount.
	Width func(b []byte) int
}

// Wrap wraps s into a paragraph of lines of length lim, with minimal
// raggedness.
func Wrap(s string, lim int) string {
	return string(WrapBytes([]byte(s), lim))
}

// WrapBytes wraps b into a paragraph of lines of length lim, with minimal
// raggedness.
func WrapBytes(b []byte, lim int) []byte {
	w := new(Wrapper)
	w.MaxWidth = lim
	return w.Wrap(b)
}

// WrapWords is superceded by the Wrapper type.
// It ignores spc.
// It should not be used in new code.
func WrapWords(words [][]byte, spc, lim, pen int) [][][]byte {
	w := new(Wrapper)
	w.MaxWidth = lim
	w.Penalty = pen
	return w.wrap(words)
}

// Wrap wraps b into a paragraph of lines of length lim, with minimal
// raggedness.
func (w *Wrapper) Wrap(b []byte) []byte {
	words := bytes.Split(bytes.Replace(bytes.TrimSpace(b), nl, sp, -1), sp)
	var lines [][]byte
	for _, line := range w.wrap(words) {
		lines = append(lines, bytes.Join(line, sp))
	}
	return bytes.Join(lines, nl)
}

func (w *Wrapper) wrap(words [][]byte) [][][]byte {
	n := len(words)
	width := w.Width
	if width == nil {
		width = utf8.RuneCount
	}
	spc := width(sp)
	wwid := make([]int, len(words))
	for i, word := range words {
		wwid[i] = width(word)
	}
	pen := w.Penalty
	if pen == 0 {
		pen = defaultPenalty * spc
	}
	lim := w.MaxWidth
	if lim == 0 {
		lim = 65 * spc
	}

	length := make([][]int, n)
	for i := 0; i < n; i++ {
		length[i] = make([]int, n)
		length[i][i] = wwid[i]
		for j := i + 1; j < n; j++ {
			length[i][j] = length[i][j-1] + spc + wwid[j]
		}
	}

	nbrk := make([]int, n)
	cost := make([]int, n)
	for i := range cost {
		cost[i] = math.MaxInt32
	}
	for i := n - 1; i >= 0; i-- {
		if length[i][n-1] <= lim || i == n-1 {
			cost[i] = 0
			nbrk[i] = n
		} else {
			for j := i + 1; j < n; j++ {
				d := lim - length[i][j-1]
				c := d*d + cost[j]
				if length[i][j-1] > lim {
					c += pen // too-long lines get a worse penalty
				}
				if c < cost[i] {
					cost[i] = c
					nbrk[i] = j
				}
			}
		}
	}

	var lines [][][]byte
	i := 0
	for i < n {
		lines = append(lines, words[i:nbrk[i]])
		i = nbrk[i]
	}
	return lines
}
