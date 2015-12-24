## Godep

[![Build Status](https://travis-ci.org/tools/godep.svg)](https://travis-ci.org/tools/godep)

[![GoDoc](https://godoc.org/github.com/tools/godep?status.svg)](https://godoc.org/github.com/tools/godep)

godep helps build packages reproducibly by fixing their dependencies.

This tool assumes you are working in a standard Go workspace, as described in
http://golang.org/doc/code.html. We expect godep to build on Go 1.4* or newer,
but you can use it on any project that works with Go 1 or newer.

## Install

```console
$ go get github.com/tools/godep
```

## How to use godep with a new project

Assuming you've got everything working already, so you can build your project
with `go install` and test it with `go test`, it's one command to start using:

```console
$ godep save -r
```

This will save a list of dependencies to the file `Godeps/Godeps.json`, copy
their source code into `Godeps/_workspace` and rewrite the dependencies. Godep
does **not copy**:

- files from source repositories that are not tracked in version control.
- `*_test.go` files.
- `testdata` directories.

Read over the contents of `Godeps/_workspace` and make sure it looks
reasonable. Then commit the whole Godeps directory to version control,
**including `Godeps/_workspace`**.

The additional flag `-r` tells save to automatically rewrite package import
paths. This allows your code to refer directly to the copied dependencies in
`Godeps/_workspace`. So, a package C that depends on package D will actually
import `C/Godeps/_workspace/src/D`. This makes C's repo self-contained and
causes `go get` to build C with the right version of all dependencies.

If you don't use `-r`, then in order to use the fixed dependencies and get
reproducible builds, you must make sure that **every time** you run a Go-related
command, you wrap it in one of these two ways:

- If the command you are running is just `go`, run it as `godep go ...`, e.g.
  `godep go install -v ./...`
- When using a different command, set your `$GOPATH` using `godep path` as
  described below.

Test files and testdata directories can be saved by adding `-t`.

## Additional Operations

### Restore

The `godep restore` command is the opposite of `godep save`. It will install the
package versions specified in `Godeps/Godeps.json` to your `$GOPATH`. This modifies the state of packages in your `$GOPATH`.

### Edit-test Cycle

1. Edit code
1. Run `godep go test`
1. (repeat)

### Add a Dependency

To add a new package foo/bar, do this:

1. Run `go get foo/bar`
1. Edit your code to import foo/bar.
1. Run `godep save` (or `godep save ./...`).

### Update a Dependency

To update a package from your `$GOPATH`, do this:

1. Run `go get -u foo/bar`
1. Run `godep update foo/bar`. (You can use the `...` wildcard, for example
`godep update foo/...`).

Before committing the change, you'll probably want to inspect the changes to
Godeps, for example with `git diff`, and make sure it looks reasonable.

## Multiple Packages

If your repository has more than one package, you're probably accustomed to
running commands like `go test ./...`, `go install ./...`, and `go fmt ./...`.
Similarly, you should run `godep save ./...` to capture the dependencies of all
packages.

## Using Other Tools

The `godep path` command helps integrate with commands other than the standard
go tool. This works with any tool that reads GOPATH from its environment, for
example the recently-released [oracle
command](http://godoc.org/code.google.com/p/go.tools/cmd/oracle).

	$ GOPATH=`godep path`:$GOPATH
	$ oracle -mode=implements .

## Old Format

Old versions of godep wrote the dependency list to a file Godeps, and didn't
copy source code. This mode no longer exists, but commands 'godep go' and 'godep
path' will continue to read the old format for some time.

## File Format

Godeps is a json file with the following structure:

```go
type Godeps struct {
	ImportPath string
	GoVersion  string   // Abridged output of 'go version'.
	Packages   []string // Arguments to godep save, if any.
	Deps       []struct {
		ImportPath string
		Comment    string // Description of commit, if present.
		Rev        string // VCS-specific commit ID.
	}
}
```

Example Godeps:

```json
{
	"ImportPath": "github.com/kr/hk",
	"GoVersion": "go1.1.2",
	"Deps": [
		{
			"ImportPath": "code.google.com/p/go-netrc/netrc",
			"Rev": "28676070ab99"
		},
		{
			"ImportPath": "github.com/kr/binarydist",
			"Rev": "3380ade90f8b0dfa3e363fd7d7e941fa857d0d13"
		}
	]
}
```

## Go 1.5 vendor/ experiment

Godep has preliminary support for the Go 1.5 vendor/
[experiment](https://github.com/golang/go/commit/183cc0cd41f06f83cb7a2490a499e3f9101befff)
utilizing the same environment variable that the go tooling itself supports:
`export GO15VENDOREXPERIMENT=1`

When `GO15VENDOREXPERIMENT=1` godep will write the vendored code into the local
package's `vendor` directory. A `Godeps/Godeps.json` file is created, just like
during normal operation. The vendor experiment is not compatible with rewrites.

There is currently no automated migration between the old Godeps workspace and
the vendor directory, but the following steps should work:

```term
$ unset GO15VENDOREXPERIMENT
$ godep restore
# The next line is only needed to automatically undo rewritten imports that were
# created with godep save -r.
$ godep save ./...
$ rm -rf Godeps
$ export GO15VENDOREXPERIMENT=1
$ godep save ./...
$ git add -A
# You should see your Godeps/_workspace/src files "moved" to vendor/.
```

NOTE: There is a "bug" in the vendor experiment that makes using `./...` with
the go tool (like go install) consider all packages inside the vendor directory:
https://github.com/golang/go/issues/11659. As a workaround you can do:

```term
$ go <cmd> $(go list ./... | grep -v /vendor/)
```

## Releasing

1. Increment the version in `version.go`.
1. Tag the commit with the same version number.
1. Update `Changelog.md`.
