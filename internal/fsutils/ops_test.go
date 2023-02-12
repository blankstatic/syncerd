package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMakeDirsForFile(t *testing.T) {
	tmpDir := t.TempDir()
	var tmp = filepath.Join(tmpDir, "sync_app", "test_dir1")

	defer os.RemoveAll(tmp)

	_, err := os.Stat(tmp)
	if err == nil {
		t.Fatalf("test dir '%v' is exist", tmp)
	}

	mderr := MakeDirsForFile(tmp + "/somefile")
	if mderr != nil {
		t.Fatal("make dirs error: ", mderr)
	}

	_, sterr := os.Stat(tmp)
	if sterr != nil {
		t.Fatalf("test dir '%v' is not exist", sterr)
	}
}

func TestGetFileSize(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.file")
	var bytes int64 = 10
	expected := fmt.Sprintf("%v bytes", bytes)

	err := CreateDummyFile(tmpFile, bytes)
	if err != nil {
		t.Fatalf("create test file %s error %v", tmpFile, err)
	}
	result := GetFileSize(tmpFile)
	if expected != result {
		t.Fatalf("get file size '%s' unexpected result '%s'", expected, result)
	}

	fictionalFile := filepath.Join(tmpDir, "fictional.file")
	result = GetFileSize(fictionalFile)
	if !strings.HasSuffix(result, "no such file or directory") {
		t.Fatalf("get file size has unexpected result '%s'", result)
	}
}

func TestCopyFileContents(t *testing.T) {
	tmpDir := t.TempDir()

	from := filepath.Join(tmpDir, "from.file")
	err := CreateDummyFile(from, 10)
	if err != nil {
		t.Fatalf("create test file %s error %v", from, err)
	}
	to := filepath.Join(tmpDir, "to.file")

	err = CopyFileContents(from, to)
	if err != nil {
		t.Fatalf("copy existing file error %v", err)
	}

	unknown := filepath.Join(tmpDir, "unknown.file")
	err = CopyFileContents(unknown, to)
	if err == nil {
		t.Fatalf("copy not existing file '%s' error", unknown)
	}
}
