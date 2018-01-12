package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blomma/viaduct/link"
	"github.com/blomma/viaduct/options"
)

// TODO: A way to exclude files, or maybe just include specific files
func main() {
	var options = options.Parse()

	// This is the path that holds the dotfiles that should be installed
	sourceDir, err := filepath.Abs(options.Path)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sourcedir: " + sourceDir)

	dotDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	targetDir, err := filepath.Abs(options.Target)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Targetdir: " + targetDir)

	if err = os.Chdir(targetDir); err != nil {
		log.Fatal(err)
	}

	if options.Unlink {
		err = filepath.Walk(sourceDir, link.Down(targetDir, sourceDir))
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// Default behaviour is to link up
	// Loop until we have succesfully linked up the dot without having to
	// unfold anything
	for {
		err = filepath.Walk(sourceDir, link.Up(targetDir, sourceDir))
		if ferr, ok := err.(*link.ErrorFoldedDirectory); ok {
			log.Println(ferr)
			dotSourceDir := filepath.Join(dotDir, ferr.Dot)

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
