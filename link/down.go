package link

import (
	"log"
	"os"
	"path/filepath"
)

// Down unlinks anything it can find in targetDir
func Down(targetDir string, sourceDir string) filepath.WalkFunc {
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

		// TODO: Check that the symlink is one we own
		fileInfo, err := os.Lstat(targetPath)

		// The only acceptable error is if the file does not exist
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if err == nil && fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			if err := os.Remove(targetPath); err != nil {
				return err
			}

			log.Println("Unlinked: " + path + " ---> " + targetPath)
		}

		return nil
	}
}
