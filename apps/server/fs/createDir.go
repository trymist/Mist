package fs

import "os"

func CreateDirIfNotExists(dirPath string, perm os.FileMode) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, perm)
		if err != nil {
			return err
		}
	}
	return nil
}
