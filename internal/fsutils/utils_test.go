package fsutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDirs(t *testing.T) {
	tmpDir := t.TempDir()
	var tmpFileDst = filepath.Join(tmpDir, "sync_app", "test_sync_dir")
	var tmpFile = filepath.Join(tmpFileDst, "some.file")
	err := MakeDirsForFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if ok := CheckDir(tmpFileDst); !ok {
		t.Fatalf("check existed dir %v error", tmpFileDst)
	}

	if err := CheckDirs(tmpFileDst); err != nil {
		t.Fatalf("check existed dirs %v error %v", tmpFileDst, err)
	}

	os.RemoveAll(tmpFileDst)

	if ok := CheckDir(tmpFileDst); ok {
		t.Fatalf("check not existed dir %v error", tmpFileDst)
	}

	if err := CheckDirs(tmpFileDst); err == nil {
		t.Fatalf("check not existed dirs %v error %v", tmpFileDst, err)
	}
}

func TestAbsPath(t *testing.T) {
	fp := ResolveAbsPath("test")
	assert.Equal(t, fp, filepath.Join(WorkingDir, "test"))

	fp = ResolveAbsPath("./")
	assert.Equal(t, fp, WorkingDir)

	fp = ResolveAbsPath("../test")
	assert.Equal(t, fp, filepath.Join(filepath.Dir(fp), "test"))
}
