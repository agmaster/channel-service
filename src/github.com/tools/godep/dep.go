package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// A Dependency is a specific revision of a package.
type Dependency struct {
	ImportPath string
	Comment    string `json:",omitempty"` // Description of commit, if present.
	Rev        string // VCS-specific commit ID.

	// used by command save & update
	ws   string // workspace
	root string // import path to repo root
	dir  string // full path to package

	// used by command update
	matched bool // selected for update by command line
	pkg     *Package

	// used by command go
	vcs *VCS
}

func eqDeps(a, b []Dependency) bool {
	ok := true
	for _, da := range a {
		for _, db := range b {
			if da.ImportPath == db.ImportPath && da.Rev != db.Rev {
				ok = false
			}
		}
	}
	return ok
}

// containsPathPrefix returns whether any string in a
// is s or a directory containing s.
// For example, pattern ["a"] matches "a" and "a/b"
// (but not "ab").
func containsPathPrefix(pats []string, s string) bool {
	for _, pat := range pats {
		if pat == s || strings.HasPrefix(s, pat+"/") {
			return true
		}
	}
	return false
}

func uniq(a []string) []string {
	var s string
	var i int
	if !sort.StringsAreSorted(a) {
		sort.Strings(a)
	}
	for _, t := range a {
		if t != s {
			a[i] = t
			i++
			s = t
		}
	}
	return a[:i]
}

// trimGoVersion and return the major version
func trimGoVersion(version string) (string, error) {
	p := strings.Split(version, ".")
	if len(p) < 2 {
		return "", fmt.Errorf("Error determing major go version from: %q", version)
	}
	return p[0] + "." + p[1], nil
}

// goVersion returns the major version string of the Go compiler
// currently installed, e.g. "go1.5".
func goVersion() (string, error) {
	// Godep might have been compiled with a different
	// version, so we can't just use runtime.Version here.
	cmd := exec.Command("go", "version")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	gv := strings.Split(string(out), " ")
	if len(gv) < 3 {
		return "", fmt.Errorf("Error splitting output of `go version`: Expected 3 or more elements, but there are < 3: %q", string(out))
	}
	return trimGoVersion(gv[2])
}
