package ratecache

import (
	"io"
	"os"
)

// DirExists checks if a directory exists.
func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// DirIsEmpty checks if the directory is empty
func DirIsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// DirIsWritable checks if the current user can access and write in the directory
func DirIsWritable(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if info.Mode().Perm()&(1<<(uint(7))) != 0 {
		return true, nil
	}
	return false, nil
}
