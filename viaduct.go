package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blomma/viaduct/flag"
	"github.com/blomma/viaduct/link"
)

// TODO: A way to exclude files, or maybe just include specific files
func main() {
	flag.Parse()

	// This is the path that holds the dotfiles that should be installed
	sourceDir, err := filepath.Abs(flag.Path)
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

	if err = os.Chdir(targetDir); err != nil {
		log.Fatal(err)
	}

	// TODO:
	if *flag.Unlink {
		err = filepath.Walk(sourceDir, link.Down(targetDir, sourceDir))
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// Default behaviour is to link up
	// Loop until we have succesfully linked up the dotdot without having to
	// unfold anything
	for {
		err = filepath.Walk(sourceDir, link.Up(targetDir, sourceDir))
		if ferr, ok := err.(*link.ErrorFoldedDirectory); ok {
			log.Println(ferr)
			dotSourceDir := filepath.Join(currentDir, ferr.Dot)

			// We need to create the actual dir that was folded
			if err = link.UnfoldAndRelink(ferr.FoldedDir, dotSourceDir, targetDir); err != nil {
				log.Fatal(err)
			}

			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		break
	}
}
