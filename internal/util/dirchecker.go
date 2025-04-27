package util

import "os"

func DirExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
