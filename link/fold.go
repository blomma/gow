package link

import (
	"log"
	"os"
	"path/filepath"
)

func UnfoldAndRelink(foldedDir string, dotSourceDir string, targetDir string) error {
	log.Println("Unlinking:" + dotSourceDir + " ---> " + targetDir)
	if err := filepath.Walk(dotSourceDir, Down(targetDir, dotSourceDir)); err != nil {
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
	return filepath.Walk(dotSourceDir, Up(targetDir, dotSourceDir))
}
