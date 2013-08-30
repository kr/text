package wrapwriter

import (
	"bytes"
	"io"
	"math"
)

var (
	nl = []byte{'\n'}
	sp = []byte{' '}
)

type Writer struct {
	width int
	buf   []byte
	bol   bool
	w     io.Writer
}

func NewWriter(w io.Writer, width int) *Writer {
	return &Writer{width: width, w: w, bol: true}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	for _, c := range p {
		w.buf = append(w.buf, c)
		if c == '\n' && w.bol {
			if err = w.wrap(); err != nil {
				return 0, err
			}
			if _, err = w.w.Write(nl); err != nil {
				return 0, err
			}
		}
		w.bol = c == '\n'
	}
	return len(p), nil
}

func (w *Writer) Flush() error {
	return w.wrap()
}

func (w *Writer) wrap() (err error) {
	b := bytes.TrimSpace(w.buf)
	w.buf = nil
	if len(b) < 1 {
		return nil
	}
	words := bytes.Split(bytes.Replace(b, nl, sp, -1), sp)
	for _, line := range wrapWords(words, w.width) {
		for i, word := range line {
			if _, err = w.w.Write(word); err != nil {
				return err
			}
			if i == len(line)-1 {
				_, err = w.w.Write(nl)
			} else {
				_, err = w.w.Write(sp)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// wrapWords is the low-level line-breaking algorithm, useful if you need more
// control over the details of the text wrapping process. For most uses, either
// Wrap or WrapBytes will be sufficient and more convenient. 
//
// wrapWords splits a list of words into lines with minimal "raggedness",
// treating each byte as one unit, accounting for 1 unit between adjacent
// words on each line, and attempting to limit lines to lim units. Raggedness
// is the total error over all lines, where error is the square of the
// difference of the length of the line and lim. Too-long lines (which only
// happen when a single word is longer than lim units) have lim**2 penalty
// units added to the error.
// TODO(kr): count runes instead of bytes for width
func wrapWords(words [][]byte, lim int) [][][]byte {
	pen := lim * lim
	n := len(words)
	length := make([][]int, n)
	for i := 0; i < n; i++ {
		length[i] = make([]int, n)
		length[i][i] = len(words[i])
		for j := i + 1; j < n; j++ {
			length[i][j] = length[i][j-1] + 1 + len(words[j])
		}
	}

	nbrk := make([]int, n)
	cost := make([]int, n)
	for i := range cost {
		cost[i] = math.MaxInt32
	}
	for i := n - 1; i >= 0; i-- {
		if length[i][n-1] <= lim {
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
