package fsutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	sep         = "/"
	defaultPerm = 0755
)

func MakeDirsForFile(path string) (err error) {
	isDir := CheckDir(path)
	var fp string
	if !isDir {
		fp = filepath.Dir(path)
		if fp == "/" || fp == "." {
			return
		}
	} else {
		fp = path
	}
	err = os.MkdirAll(fp, defaultPerm)
	return
}

func GetFileSize(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%v bytes", info.Size())
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func GetSubdirs(root string) (dirs []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return
}
