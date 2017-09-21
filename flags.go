package main

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
	flagVersion    = flag.Bool("v", false, "Show the version number")
	flagVersionAll = flag.Bool("V", false, "Show full version information")
	flagUnlink     = flag.Bool("u", false, "Unlink")
	flagPath       = ""
)

func commandLineFlags() {
	flag.Parse()
	if *flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *flagVersionAll {
		fmt.Println("Version:", Version)
		fmt.Println("BuildNumber:", BuildNumber)
		fmt.Println("CommitHash:", CommitHash)
		fmt.Println("CompileDate:", CompileDate)
		os.Exit(0)
	}

	flagPath = flag.Arg(0)

	if flagPath == "" {
		fmt.Println("Missing argument for path")
		os.Exit(0)
	}
}
