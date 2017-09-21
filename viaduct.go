package main

import (
	"log"
	"os"
	"path/filepath"
)

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

func existsAndSymlink(symLinkPath string) bool {
	fileInfo, err := os.Lstat(symLinkPath)
	if err != nil {
		return false
	}

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}

	return false
}

func linkDown(targetDir string, sourceDir string) filepath.WalkFunc {
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
		if err != nil {
			log.Fatal(err)
		}

		if existsAndSymlink(targetPath) {
			err := os.Remove(targetPath)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Unlinked: " + targetPath)
		}

		return nil
	}
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
	sourceDir, err := filepath.Abs(flagPath)
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

	if *flagUnlink {
		err = filepath.Walk(sourceDir, linkDown(targetDir, sourceDir))
	} else {
		err = filepath.Walk(sourceDir, linkUp(targetDir, sourceDir))
	}

	if err != nil {
		log.Fatal(err)
	}
}
