package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// exists returns whether the given file or directory exists or not
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
			log.Print(err)
			return nil
		}

		relSourcePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			log.Fatal(err)
		}

		if relSourcePath == "." {
			return nil
		}

		targetPath := filepath.Join(targetDir, relSourcePath)

		log.Println(targetPath + " ---> " + relSourcePath)
		if !exists(targetPath) {
			log.Println("Linking " + targetPath + " ---> " + path)
			// err := os.Symlink(path, targetPath)
			// if err != nil {
			// 	log.Fatal(err)
			// }
		}

		return nil
	}
}

func main() {
	sourceDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	sourceDir = filepath.Clean(sourceDir)
	fmt.Println("Sourcedir --> " + sourceDir)

	pwd, err := os.Getwd()
	targetDir := filepath.Join(pwd, "..")
	targetDir = filepath.Clean(targetDir)
	fmt.Println("Targetdir --> " + targetDir)

	err = filepath.Walk(sourceDir, linkUp(targetDir, sourceDir))
	if err != nil {
		log.Fatal(err)
	}
}
