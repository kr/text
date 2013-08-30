// Command tab translates tabs to spaces.
//
//   Usage: tab [file...]
//
// Uses the elastic tabs algorithm. See tabwriter and original
// source.
package main

import (
	"io"
	"log"
	"os"
	"text/tabwriter"
)

var (
	minw      = 2
	tabw      = 4
	npad      = 2
	padc byte = ' '
	flag uint = 0
)

func main() {
	args := os.Args[1:]
	w := tabwriter.NewWriter(os.Stdout, minw, tabw, npad, padc, flag)
	if len(args) > 0 {
		for _, s := range args {
			if f, err := os.Open(s); err == nil {
				copyin(w, f)
			} else {
				log.Println(err)
			}
		}
	} else {
		copyin(w, os.Stdin)
	}
}

func copyin(w *tabwriter.Writer, r io.Reader) {
	if _, err := io.Copy(w, r); err != nil {
		log.Println(err)
	}
	if err := w.Flush(); err != nil {
		log.Println(err)
	}
}
