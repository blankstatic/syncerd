package syncd

import (
	"context"
	"path/filepath"
	"syncer/internal/fsutils"
	"syncer/internal/logging"
	"testing"
	"time"
)

func init() {
	logging.SetupForTests()
}

func TestBaseSync(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	dst := filepath.Join(tmp, "dst")
	err := fsutils.MakeDirsForFile(filepath.Join(src, "something"))
	if err != nil {
		t.Fatalf("create dir '%s' error %v", dst, err)
	}
	err = fsutils.MakeDirsForFile(filepath.Join(dst, "something"))
	if err != nil {
		t.Fatalf("create dir '%s' error %v", dst, err)
	}

	bs := NewBaseSync(ctx, src, dst, Options{Force: true, Interval: time.Minute * 1})
	func() {
		<-time.After(time.Millisecond * 10)
		cancel()
	}()

	srcFiles := map[string]int64{"skip.file": 10, "test.file": 10, "mod.file": 10}
	for file, size := range srcFiles {
		err := fsutils.CreateDummyFile(filepath.Join(src, file), size)
		if err != nil {
			t.Fatalf("create file %s error %v", file, err)
		}
	}

	dstFiles := map[string]int64{"skip.file": 10, "mod.file": 11, "removed.file": 10}
	for file, size := range dstFiles {
		err := fsutils.CreateDummyFile(filepath.Join(dst, file), size)
		if err != nil {
			t.Fatalf("create file %s error %v", file, err)
		}
	}

	bs.Startup()
	bs.RunSyncLoop()
	bs.Cleanup()

	bs.RunSyncLoop() // empty actions
}
