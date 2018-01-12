package options

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

var (
	version    = flag.Bool("v", false, "Show the version number")
	versionAll = flag.Bool("V", false, "Show full version information")
	unlink     = flag.Bool("u", false, "Unlink")
	target     = flag.String("t", "..", "Targetpath, default is directory above .dotfiles")
	path       = ""
)

// Options holds options passed to the program
type Options struct {
	Version    bool
	VersionAll bool
	Unlink     bool
	Target     string
	Path       string
}

// Parse parses options passed to the program
func Parse() Options {
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

	path = flag.Arg(0)

	if path == "" {
		fmt.Println("Missing argument for path")
		os.Exit(0)
	}

	return Options{
		Path:       path,
		Target:     *target,
		Unlink:     *unlink,
		Version:    *version,
		VersionAll: *versionAll}
}
