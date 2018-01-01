package link

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func getRelPathInsideSource(targetPath string, sourceDir string) (string, error) {
	originPath, err := os.Readlink(targetPath)
	if err != nil {
		return "", err
	}

	s := strings.LastIndex(sourceDir, string(os.PathSeparator))
	relSourcePath, err := filepath.Rel(sourceDir[0:s], originPath)
	if err != nil {
		log.Print(err)
		return "", err
	}

	if strings.Contains(relSourcePath, "../") {
		return "", &ErrorNotOwned{Message: "Path: " + targetPath + "is not inside sourcedir:" + sourceDir + ", it points to: " + originPath}
	}

	return relSourcePath, nil
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
	relSourcePath, err := getRelPathInsideSource(targetPath, sourceDir)
	if err != nil {
		return err
	}

	parts := strings.Split(relSourcePath, string(os.PathSeparator))
	return &ErrorFoldedDirectory{Message: "Found folded directory:" + targetPath, FoldedDir: targetPath, Dot: parts[0]}
}

func Up(targetDir string, sourceDir string) filepath.WalkFunc {
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
		exists, err := exists(targetPath)
		if err != nil {
			return err
		}
		if !exists {
			if err := os.Symlink(path, targetPath); err != nil {
				return err
			}
			log.Println("Linked: " + path + " ---> " + targetPath)
		} else {
			if err := isFolded(targetPath, sourceDir); err != nil {
				return err
			}
			log.Println("Exists: " + path + " ---> " + targetPath)
		}

		return nil
	}
}
