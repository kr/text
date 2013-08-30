// Command flow reflows text in paragraphs.
//
//   Usage: flow [-w width] [-pre prefix] [file...]
package main

import (
	"flag"
	"github.com/kr/pty"
	"github.com/kr/text/flowwriter"
	"io"
	"log"
	"os"
)

func main() {
	var width int
	var pre string
	flag.IntVar(&width, "w", 0, "width")
	flag.StringVar(&pre, "pre", "    ", "preformatted prefix")
	flag.Parse()
	if width < 1 {
		_, width, _ = pty.Getsize(os.Stdout)
	}
	if width < 1 {
		width = 80
	}

	w := flowwriter.NewWriter(os.Stdout, width)
	w.SetNoWrap(pre)
	if flag.NArg() > 0 {
		for _, s := range flag.Args() {
			if f, err := os.Open(s); err == nil {
				copyin(w, f)
				f.Close()
			} else {
				log.Println(err)
			}
		}
	} else {
		copyin(w, os.Stdin)
	}
}

func copyin(w *flowwriter.Writer, r io.Reader) {
	if _, err := io.Copy(w, r); err != nil {
		log.Println(err)
	}
	if err := w.Flush(); err != nil {
		log.Println(err)
	}
}
