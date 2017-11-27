package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
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

func existsAndSymlink(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false
	}

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}

	return false
}

func isFolded(targetPath string, sourceDir string) error {
	fileInfo, err := os.Lstat(targetPath)
	if err != nil {
		return err
	}

	if fileInfo.Mode()&os.ModeSymlink != os.ModeSymlink {
		return nil
	}

	// Check if we own this
	originPath, err := os.Readlink(targetPath)
	if err != nil {
		return err
	}

	// sourcedir is always the top dir in the dotfiles folder, we strip of the last part and check
	s := strings.LastIndex(sourceDir, string(os.PathSeparator))
	relSourcePath, err := filepath.Rel(sourceDir[0:s], originPath)
	if err != nil {
		log.Print(err)
		return err
	}

	if strings.Contains(relSourcePath, "../") {
		return nil
	}

	parts := strings.Split(relSourcePath, string(os.PathSeparator))
	return &errorFoldedDirectory{message: "Found folded directory:" + targetPath, foldedDir: targetPath, dot: parts[0]}
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
		if existsAndSymlink(targetPath) {
			if err := os.Remove(targetPath); err != nil {
				log.Fatal(err)
			}
			log.Println("Unlinked: " + path + " ---> " + targetPath)
		}

		return nil
	}
}

type errorFoldedDirectory struct {
	message   string
	dot       string
	foldedDir string
}

func (e *errorFoldedDirectory) Error() string {
	return e.message
}

func linkUp(targetDir string, sourceDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relSourcePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if relSourcePath == "." {
			return nil
		}

		targetPath := filepath.Join(targetDir, relSourcePath)
		if !exists(targetPath) {
			if err := os.Symlink(path, targetPath); err != nil {
				return err
			}
			log.Println("Linked: " + path + " ---> " + targetPath)
		} else {
			if ferr := isFolded(targetPath, sourceDir); ferr != nil {
				return ferr
			}
			log.Println("Exists: " + path + " ---> " + targetPath)
		}

		return nil
	}
}

func unfoldAndRelink(foldedDir string, dotSourceDir string, targetDir string) error {
	log.Println("Unlinking:" + dotSourceDir + " ---> " + targetDir)
	if err := filepath.Walk(dotSourceDir, linkDown(targetDir, dotSourceDir)); err != nil {
		return err
	}

	// Create with same perm as parent dir
	parentDir := filepath.Join(foldedDir, "..")
	fileInfo, err := os.Lstat(parentDir)
	if err != nil {
		return err
	}

	if err = os.Mkdir(foldedDir, fileInfo.Mode()); err != nil {
		return err
	}

	log.Println("Relinking:" + dotSourceDir + " ---> " + targetDir)
	return filepath.Walk(dotSourceDir, linkUp(targetDir, dotSourceDir))
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

	if err = os.Chdir(targetDir); err != nil {
		log.Fatal(err)
	}

	if *flagUnlink {
		err = filepath.Walk(sourceDir, linkDown(targetDir, sourceDir))
	} else {
		for {
			err = filepath.Walk(sourceDir, linkUp(targetDir, sourceDir))
			if ferr, ok := err.(*errorFoldedDirectory); ok {
				log.Println(ferr)
				dotSourceDir := filepath.Join(currentDir, ferr.dot)

				// We need to create the actual dir that was folded
				if err = unfoldAndRelink(ferr.foldedDir, dotSourceDir, targetDir); err != nil {
					log.Fatal(err)
				}
			} else {
				break
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}
