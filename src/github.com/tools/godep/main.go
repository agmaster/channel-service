package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"text/template"
)

var (
	cpuprofile       string
	verbose          bool // Verbose flag for commands that support it
	debug            bool // Debug flag for commands that support it
	majorGoVersion   string
	VendorExperiment bool
	sep              string
)

// Command is an implementation of a godep command
// like godep save or godep go.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// Name of the command
	Name string

	// Args the command would expect
	Args string

	// Short is the short description shown in the 'godep help' output.
	Short string

	// Long is the long message shown in the
	// 'godep help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

// UsageExit prints usage information and exits.
func (c *Command) UsageExit() {
	fmt.Fprintf(os.Stderr, "Args: godep %s [-v] [-d] %s\n\n", c.Name, c.Args)
	fmt.Fprintf(os.Stderr, "Run 'godep help %s' for help.\n", c.Name)
	os.Exit(2)
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed
// by 'godep help'.
var commands = []*Command{
	cmdSave,
	cmdGo,
	cmdGet,
	cmdPath,
	cmdRestore,
	cmdUpdate,
	cmdDiff,
	cmdVersion,
}

func main() {
	flag.Usage = usageExit
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("godep: ")
	args := flag.Args()
	if len(args) < 1 {
		usageExit()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	var err error
	majorGoVersion, err = goVersion()
	if err != nil {
		log.Fatal(err)
	}

	// VendorExperiment is the Go 1.5 vendor directory experiment flag, see
	// https://github.com/golang/go/commit/183cc0cd41f06f83cb7a2490a499e3f9101befff
	VendorExperiment = (majorGoVersion == "go1.5" && os.Getenv("GO15VENDOREXPERIMENT") == "1") ||
		(majorGoVersion == "go1.6" && os.Getenv("GO15VENDOREXPERIMENT") != "0")

	// sep is the signature set of path elements that
	// precede the original path of an imported package.
	sep = defaultSep(VendorExperiment)

	for _, cmd := range commands {
		if cmd.Name == args[0] {
			cmd.Flag.BoolVar(&verbose, "v", false, "enable verbose output")
			cmd.Flag.BoolVar(&debug, "d", false, "enable debug output")
			cmd.Flag.StringVar(&cpuprofile, "cpuprofile", "", "Write cpu profile to this file")
			cmd.Flag.Usage = func() { cmd.UsageExit() }
			cmd.Flag.Parse(args[1:])

			debugln("majorGoVersion", majorGoVersion)
			debugln("VendorExperiment", VendorExperiment)
			debugln("sep", sep)

			if cpuprofile != "" {
				f, err := os.Create(cpuprofile)
				if err != nil {
					log.Fatal(err)
				}
				pprof.StartCPUProfile(f)
				defer pprof.StopCPUProfile()
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "godep: unknown command %q\n", args[0])
	fmt.Fprintf(os.Stderr, "Run 'godep help' for usage.\n")
	os.Exit(2)
}

var usageTemplate = `
Godep is a tool for managing Go package dependencies.

Usage:

	godep command [arguments]

The commands are:
{{range .}}
    {{.Name | printf "%-8s"}} {{.Short}}{{end}}

Use "godep help [command]" for more information about a command.
`

var helpTemplate = `
Args: godep {{.Name}} [-v] [-d] {{.Args}}

{{.Long | trim}}

If -v is given, verbose output is enabled.

If -d is given, debug output is enabled (you probably don't want this, see -v).

`

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: godep help command\n\n")
		fmt.Fprintf(os.Stderr, "Too many arguments given.\n")
		os.Exit(2)
	}
	for _, cmd := range commands {
		if cmd.Name == args[0] {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}
}

func usageExit() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
	})
	template.Must(t.Parse(strings.TrimSpace(text) + "\n\n"))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
