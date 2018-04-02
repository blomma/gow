package option

import (
	"flag"
	"fmt"
	"os"
)

// These variables should be set by the linker when compiling
var (
	Version     = "0.0.0-unknown"
	BuildNumber = "unknown"
	CommitHash  = "Unknown"
	CompileDate = "Unknown"
)

// Options holds options passed to the program
type Options struct {
	Unlink bool
	Target string
	Path   string
}

// Parse parses options passed to the program
func (o *Options) Parse() {
	var version = flag.Bool("v", false, "Show the version number")
	var versionAll = flag.Bool("V", false, "Show the full version information")

	flag.BoolVar(&o.Unlink, "u", false, "Unlink")
	flag.StringVar(&o.Target, "t", "..", "Targetpath, default is directory above .dotfiles")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *versionAll {
		fmt.Println("Version:", Version)
		fmt.Println("BuildNumber:", BuildNumber)
		fmt.Println("CommitHash:", CommitHash)
		fmt.Println("CompileDate:", CompileDate)
		os.Exit(0)
	}

	o.Path = flag.Arg(0)

	if o.Path == "" {
		fmt.Println("Missing argument for path")
		os.Exit(0)
	}
}
