package flag

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
	Unlink     = flag.Bool("u", false, "Unlink")
	Path       = ""
)

func Parse() {
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

	Path = flag.Arg(0)

	if Path == "" {
		fmt.Println("Missing argument for path")
		os.Exit(0)
	}
}
