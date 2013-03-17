package colwriter

import (
	"bytes"
	"testing"
)

var s = `
.git
.gitignore
.godir
Procfile:
README.md
api.go
apps.go
auth.go
darwin.go
data.go
dyno.go:
env.go
git.go
help.go
hkdist
linux.go
ls.go
main.go
plugin.go
run.go
scale.go
ssh.go
tail.go
term
unix.go
update.go
version.go
windows.go
`[1:]

var de = `
.git       README.md  darwin.go  git.go     ls.go      scale.go   unix.go
.gitignore api.go     data.go    help.go    main.go    ssh.go     update.go
.godir     apps.go    dyno.go:   hkdist     plugin.go  tail.go    version.go
Procfile:  auth.go    env.go     linux.go   run.go     term       windows.go
`[1:]

var ce = `
.git       .gitignore .godir

Procfile:
README.md api.go    apps.go   auth.go   darwin.go data.go

dyno.go:
env.go     hkdist     main.go    scale.go   term       version.go
git.go     linux.go   plugin.go  ssh.go     unix.go    windows.go
help.go    ls.go      run.go     tail.go    update.go
`[1:]

func TestColumns(t *testing.T) {
	b := new(bytes.Buffer)
	w := NewWriter(b, 80, 0)
	if _, err := w.Write([]byte(s)); err != nil {
		t.Error(err)
	}
	if err := w.Flush(); err != nil {
		t.Error(err)
	}
	g := string(b.Bytes())
	if de != g {
		t.Log("\n" + de)
		t.Log("\n" + g)
		t.Errorf("%q != %q", de, g)
	}
}

func TestColumnsColon(t *testing.T) {
	b := new(bytes.Buffer)
	w := NewWriter(b, 80, BreakOnColon)
	if _, err := w.Write([]byte(s)); err != nil {
		t.Error(err)
	}
	if err := w.Flush(); err != nil {
		t.Error(err)
	}
	g := string(b.Bytes())
	if ce != g {
		t.Log("\n" + ce)
		t.Log("\n" + g)
		t.Errorf("%q != %q", ce, g)
	}
}

func TestColWriterFlushEmpty(t *testing.T) {
	b := new(bytes.Buffer)
	w := NewWriter(b, 80, 0)
	if err := w.Flush(); err != nil {
		t.Error(err)
	}
	if g := b.Bytes(); len(g) != 0 {
		t.Errorf("expected empty output, got %#v", g)
	}
}
