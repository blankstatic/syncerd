package fsutils

import (
	"path/filepath"
	"testing"
)

func TestMD5Files(t *testing.T) {
	tmpDir := t.TempDir()
	tmp1Name := "tmp1.file"
	tmp2Name := "tmp2.file"
	var tmp1 = filepath.Join(tmpDir, tmp1Name)
	var tmp2 = filepath.Join(tmpDir, tmp2Name)
	err := CreateDummyFile(tmp1, 10)
	if err != nil {
		t.Fatalf("create test file %s error %v", tmp1, err)
	}
	err = CreateDummyFile(tmp2, 10)
	if err != nil {
		t.Fatalf("create test file %s error %v", tmp2, err)
	}

	files, err := MD5All(tmpDir, false)
	if err != nil {
		t.Fatalf("get md5 files error %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("get md5 files error %v", files)
	}

	hash1, exist := files[tmp1]
	if !exist {
		t.Fatalf("md5 file %s not exist", tmp1Name)
	}
	hash2, exist := files[tmp2]
	if !exist {
		t.Fatalf("md5 file %s not exist", tmp2Name)
	}
	if hash1 != hash2 {
		t.Fatalf("md5 files same content has diff hash %s != %s", hash1, hash2)
	}

	// remove root from path
	files, err = MD5All(tmpDir, true)
	if err != nil {
		t.Fatalf("get md5 files error %v", err)
	}
	_, exist = files[sep+tmp1Name]
	if !exist {
		t.Fatalf("md5 file %s not exist", tmp1Name)
	}
}
