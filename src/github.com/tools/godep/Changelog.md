# v44 2015/12/23

* Clean package roots when attempting to find a vendor directory so we don't loop forever.
    * Fixes 382
  
# v43 2015/12/22

* Better error messages when parsing Godeps.json: Fixes #372

# v42 2015/12/22

* Fix a bunch of GO15VENDOREXPERIMENT issues
    * Find package directories better. Previously we used build.FindOnly which didn't work the way I expected it to (any dir would work w/o error).
    * Set the VendorExperiment bool based on go version as 1.6 defaults to on.
    * A bunch of extra debugging for use while sanity checking myself.
    * vendor flag for test structs.
    * Some tests for vendor/ stuff:
        * Basic Test
        * Transitive
        * Transitive, across GOPATHs + collapse vendor/ directories.
* Should Fix #358

# v41 2015/12/17

* Don't rewrite packages outside of the project. This would happen if you specified
  an external package for vendoring when you ran `goodep save -r ./... github.com/some/other/package`

# v40 2015/12/17

* When downloading a dependency, create the base directory if needed.

# v39 2015/12/16

* Record only the major go version (ex. go1.5) instead of the complete string.

# v38 2015/12/16

* Replace `go get`, further fix up restore error handling/reporting.
    * Fixes #186
    * Don't bother restoring/downloading if already done.

# v37 2015/12/15

* Change up how download/restore works a little
    * Try to load the package after downloading/restoring. Previously
      that was done too early in the process.
    * make previous verbose output debug output
    * report a typed error instead of a string from listPackage so it can
      be asserted to provide a nicer error.
    * Catch go get errors that say there are no go files found. See code
      comment as to why.
    * do *all* downloading during download phase.

# v36 2015/12/14

* Fixes #358: Using wrong variable. Will add test after release.

# v35 2015/12/11

* Fixes #356: Major performance regressions in v34
    * Enable cpu profiling via flag on save.
    * Cache packages by dir
    * Don't do a full import pass on deps for packages in the GOROOT
    * create a bit less garbage at times
* Generalize -v & -d flags

# v34 2015/12/08

* We now use build.Context to help locate packages only and do our own parsing (via go/ast).
* Fixes reported issues caused by v33 (Removal of `go list`):
    * #345: Bug in godep restore
    * #346: Fix loading a dot package
    * #348: Godep save issue when importing lib/pq
    * #350: undefined: build.MultiplePackageError
    * #351: stow away helper files
    * #353: cannot find package "appengine"

# v33 2015/12/07

* Replace the use of `go list`. This is a large change although all existing tests pass.

# v32 2015/12/02

* Eval Symlinks in Contains() check.

# v31 2015/12/02

* In restore, mention which package had the problem -- @shurcool

# v30 2015/11/25

* Add `-t` flag to the `godep get` command.

# v29 2015/11/17

* Temp work around to fix issue with LICENSE files.

# v28 2015/11/09

* Make `version` an actual command.

# v27 2015/11/06

* run command once during restore -v

# v26 2015/11/05

* Better fix for the issue fixed in v25: All update paths are now path.Clean()'d

# v25 2015/11/05

* `godep update package/` == `godep update package`. Fixes #313

# v24 2015/11/05

* Honor -t during update. Fixes #312

# v23 2015/11/05

* Do not use --debug to find full revision name for mercurial repositories

# v22 2016/11/14

* s/GOVENDOREXPERIMENT/GO15VENDOREXPERIMENT :-(

# v21 2016/11/13

* Fix #310: Case insensitive fs issue

# v20 2016/11/13

* Attempt to include license files when vendoring. (@client9)

# v19 2016/11/3

* Fix conflict error message. Revisions were swapped. Also better selection of package that needs update.

# v18 2016/10/16

* Improve error message when trying to save a conflicting revision.

# v17 2016/10/15

* Fix for v16 bug. All vcs list commands now produce paths relative to the root of the vcs.

# v16 2015/10/15

* Determine repo root using vcs commands and use that instead of dep.dir

# v15 2015/10/14

* Update .travis.yml file to do releases to github

# v14 2015/10/08

* Don't print out a workspace path when GO15VENDOREXPERIMENT is active. The vendor/ directory is not a valid workspace, so can't be added to your $GOPATH.

# v13 2015/10/07

* Do restores in 2 separate steps, first download all deps and then check out the recorded revisions.
* Update Changelog date format

# v12 2015/09/22

* Extract errors into separate file.

# v11 2015/08/22

* Amend code to pass golint.

# v10 2015/09/21

* Analyse vendored package test dependencies.
* Update documentation.

# v9 2015/09/17

* Don't save test dependencies by default.

# v8 2015/09/17

* Reorganize code.

# v7 2015/09/09

* Add verbose flag.
* Skip untracked files.
* Add VCS list command.

# v6 2015/09/04

*  Revert ignoring testdata directories and instead ignore it while
processing Go files and copy the whole directory unmodified.

# v5 2015/09/04

* Fix vcs selection in restore command to work as go get does

# v4 2015/09/03

* Remove the deprecated copy option.

# v3 2015/08/26

* Ignore testdata directories

# v2 2015/08/11

* Include command line packages in the set to copy

This is a simplification to how we define the behavior
of the save command. Now it has two distinct package
parameters, the "root set" and the "destination", and
they have clearer roles. The packages listed on the
command line form the root set; they and all their
dependencies will be copied into the Godeps directory.
Additionally, the destination (always ".") will form the
initial list of "seen" import paths to exclude from
copying.

In the common case, the root set is equal to the
destination, so the effective behavior doesn't change.
This is primarily just a simpler definition. However, if
the user specifies a package on the command line that
lives outside of . then that package will be copied.

As a side effect, there's a simplification to the way we
add packages to the initial "seen" set. Formerly, to
avoid copying dependencies unnecessarily, we would try
to find the root of the VCS repo for each package in the
root set, and mark the import path of the entire repo as
seen. This meant for a repo at path C, if destination
C/S imports C/T, we would not copy C/T into C/S/Godeps.
Now we don't treat the repo root specially, and as
mentioned above, the destination alone is considered
seen.

This also means we don't require listed packages to be
in VCS unless they're outside of the destination.

# v1 2015/07/20

* godep version command

Output the version as well as some godep runtime information that is
useful for debugging user's issues.

The version const would be bumped each time a PR is merged into master
to ensure that we'll be able to tell which version someone got when they
did a `go get github.com/tools/godep`.

# Older changes

Many and more, see `git log -p`
