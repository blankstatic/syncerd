package lck

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLockCycle(t *testing.T) {
	defer func(fn func() string) {
		GetLockFilename = fn
	}(GetLockFilename)

	origLockFile := GetLockFilename()
	if !strings.Contains(origLockFile, lockSrc) && !strings.HasSuffix(origLockFile, lockSuffix) {
		t.Fatalf("lock filename has wrong value %s", origLockFile)
	}

	tmpLock := filepath.Join(t.TempDir(), "tmp.lock")
	GetLockFilename = func() string {
		return tmpLock
	}

	if err := Lock(); err != nil {
		t.Fatalf("unexpected lock error: %v", err)
	}

	lckd, lockFile := IsLocked()
	if !lckd {
		t.Fatalf("lock '%s' failed", lockFile)
	}

	if err := Unlock(); err != nil {
		t.Fatalf("unexpected unlock error: %v", err)
	}

	lckd, lockFile = IsLocked()
	if lckd {
		t.Fatalf("unlock '%s' failed", lockFile)
	}
}
