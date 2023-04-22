package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
)

var WorkingDir, _ = os.Getwd()

func CheckDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func CheckDirs(dirs ...string) error {
	for _, dir := range dirs {
		ok := CheckDir(dir)
		if !ok {
			return fmt.Errorf("check dir %s failed", dir)
		}
	}
	return nil
}

// tests used
func CreateDummyFile(filename string, size int64) error {
	err := MakeDirsForFile(filename)
	if err != nil {
		return err
	}
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = fd.Seek(size-1, 0)
	if err != nil {
		return err
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		return err
	}
	err = fd.Close()
	return err
}

func ResolveAbsPath(path string) string {
	fp := filepath.Clean(path)
	if isAbs := filepath.IsAbs(fp); isAbs {
		return fp
	}
	return filepath.Join(WorkingDir, fp)
}
