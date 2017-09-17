package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func linkUp(targetDir string, sourceDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		relSourcePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			log.Fatal(err)
		}

		if relSourcePath == "." {
			return nil
		}

		targetPath := filepath.Join(targetDir, relSourcePath)
		relativePath, err := filepath.Rel(targetDir, path)
		if err != nil {
			log.Fatal(err)
		}

		if !exists(targetPath) {
			err := os.Symlink(relativePath, targetPath)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Linked: " + relativePath + " ---> " + targetPath)
		} else {
			log.Println("Exists: " + relativePath + " ---> " + targetPath)
		}

		return nil
	}
}

// TODO: A way to exclude files, or maybe just include specific files
func main() {
	commandLineFlags()

	// This is the path that holds the dotfiles that should be installed
	sourceDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sourcedir: " + sourceDir)

	// This is where we should install the files from sourceDir, it is
	// hardcoded to the dir above the current dir
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	targetDir := filepath.Join(currentDir, "..")
	log.Println("Targetdir: " + targetDir)

	err = os.Chdir(targetDir)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(sourceDir, linkUp(targetDir, sourceDir))
	if err != nil {
		log.Fatal(err)
	}
}
