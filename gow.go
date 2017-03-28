package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func linkUp(targetDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		// If this is a directory check if it exists
		if info.IsDir() {
			result, _ := exists(path)
			if !result {
				// Create it as a symlink
			}

			return nil
		}

		return nil
	}
}

func main() {
	// Directory to symlink
	dir := os.Args[1]
	fmt.Println(dir)

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ex)

	// Current dir
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	parentDir := filepath.Join(exPath, "..")
	fmt.Println(parentDir)

	err = filepath.Walk(dir, linkUp(parentDir))
	if err != nil {
		log.Fatal(err)
	}

	// target := "symtarget.txt"
	// // os.MkdirAll(path, 0755)
	// ioutil.WriteFile(target, []byte("Hello\n"), 0644)
	// symlink := "symlink"
	// error := os.Symlink(target, symlink)
	// fmt.Println(error)
}
