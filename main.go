package main

import (
	"log"
	"path/filepath"

	"github.com/blomma/viaduct/link"
	"github.com/blomma/viaduct/option"
)

// TODO: A way to exclude files, or maybe just include specific files
func main() {
	var options = option.Options{}
	options.Parse()

	// This is the path that holds the dotfiles that should be installed
	dotSourceDir, err := filepath.Abs(options.Path)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sourcedir: " + dotSourceDir)

	targetDir, err := filepath.Abs(options.Target)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Targetdir: " + targetDir)

	if options.Unlink {
		err = filepath.Walk(dotSourceDir, link.Down(targetDir, dotSourceDir))
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// Default behavior is to link up
	// Loop until we have successfully linked up the dot without having to
	// unfold anything
	for {
		err = filepath.Walk(dotSourceDir, link.Up(targetDir, dotSourceDir))
		if ferr, ok := err.(*link.ErrorFoldedDirectory); ok {
			dotDir := filepath.Join(dotSourceDir, "..")
			foldedDotSourceDir := filepath.Join(dotDir, ferr.Dot)

			// We need to create the actual dir that was folded
			if err = link.UnfoldAndRelink(ferr.FoldedDir, foldedDotSourceDir, targetDir); err != nil {
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
