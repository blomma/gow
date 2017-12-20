package link

import "os"

func existsAndSymlink(path string) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true, nil
	}

	return false, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return false, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
