package link

import (
	"log"
	"os"
	"path/filepath"
)

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
		existsAndSymlink, err := existsAndSymlink(targetPath)
		if err != nil {
			return err
		}
		if existsAndSymlink {
			if err := os.Remove(targetPath); err != nil {
				return err
			}
			log.Println("Unlinked: " + path + " ---> " + targetPath)
		}

		return nil
	}
}
