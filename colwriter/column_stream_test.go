package colwriter_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/kr/text/colwriter"
)

func writeChunks(t *testing.T, width int, flag uint, chunks []string) string {
	t.Helper()
	var output bytes.Buffer
	w := colwriter.NewWriter(&output, width, flag)
	for _, chunk := range chunks {
		if n, err := w.Write([]byte(chunk)); n != len(chunk) || err != nil {
			t.Fatalf("Write(%q) = (%d, %v), want (%d, nil)", chunk, n, err, len(chunk))
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
	return output.String()
}

func TestWriterChunkedInput(t *testing.T) {
	for _, test := range []struct {
		name, input, want string
		width             int
		flag              uint
	}{
		{"tab in split line", "a\tb\n", "a   b\n", 20, 0},
		{"colon before split newline", "heading:\na\nb\n", "heading:\na b\n", 20, colwriter.BreakOnColon},
		{"ordinary lines", "a\nb\n", "a b\n", 20, 0},
		{"colon flag disabled", "heading:\na\nb\n", "heading: b\na\n", 20, 0},
		{"tab at line start", "\tb\n", "    b\n", 20, 0},
		{"newline resets column", "abc\n\tb\n", "abc\n    b\n", 1, 0},
		{"empty input", "", "", 20, colwriter.BreakOnColon},
	} {
		t.Run(test.name, func(t *testing.T) {
			type partition struct {
				name   string
				chunks []string
			}
			partitions := []partition{{"single", []string{test.input}}}
			for split := 0; split <= len(test.input); split++ {
				partitions = append(partitions, partition{
					fmt.Sprintf("split-%d", split), []string{test.input[:split], test.input[split:]},
				})
			}
			var bytewise []string
			emptyWrites := []string{""}
			for i := 0; i < len(test.input); i++ {
				bytewise = append(bytewise, test.input[i:i+1])
				emptyWrites = append(emptyWrites, test.input[i:i+1], "")
			}
			partitions = append(partitions,
				partition{"bytewise", bytewise},
				partition{"empty-interleaved", emptyWrites},
				partition{"complete-lines", strings.SplitAfter(test.input, "\n")},
			)
			for _, partition := range partitions {
				t.Run(partition.name, func(t *testing.T) {
					if got := writeChunks(t, test.width, test.flag, partition.chunks); got != test.want {
						t.Errorf("output = %q, want %q", got, test.want)
					}
				})
			}
		})
	}
}

func TestWriterFlushReuse(t *testing.T) {
	for _, test := range []struct {
		name, first, second, want string
		flag                      uint
	}{
		{"tab state", "a", "\tb\n", "a\n    b\n", 0},
		{"colon state", "heading:", "\na\n", "heading:\n  a\n", colwriter.BreakOnColon},
	} {
		t.Run(test.name, func(t *testing.T) {
			fresh := writeChunks(t, 20, test.flag, []string{test.first}) +
				writeChunks(t, 20, test.flag, []string{test.second})
			var output bytes.Buffer
			w := colwriter.NewWriter(&output, 20, test.flag)
			for _, batch := range []string{test.first, test.second} {
				if n, err := w.Write([]byte(batch)); n != len(batch) || err != nil {
					t.Fatalf("Write(%q) = (%d, %v), want (%d, nil)", batch, n, err, len(batch))
				}
				if err := w.Flush(); err != nil {
					t.Fatal(err)
				}
			}
			if got := output.String(); got != test.want || fresh != test.want {
				t.Errorf("reused output = %q, fresh output = %q, want %q", got, fresh, test.want)
			}
		})
	}
}
