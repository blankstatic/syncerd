package lck

import (
	"os"
	"path/filepath"
)

const (
	lockSrc    = "/tmp/"
	lockSuffix = ".lock"
)

var GetLockFilename = func() string {
	return lockSrc + filepath.Base(os.Args[0]) + lockSuffix
}

func IsLocked() (bool, string) {
	lck := GetLockFilename()
	_, err := os.Stat(lck)
	if err == nil {
		return true, lck
	}
	Lock()
	return false, lck
}

func Lock() error {
	f, err := os.Create(GetLockFilename())
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func Unlock() error {
	err := os.Remove(GetLockFilename())
	return err
}
