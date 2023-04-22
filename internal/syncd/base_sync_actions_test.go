package syncd

import (
	"os"
	"path/filepath"
	"syncer/internal/fsutils"
	"syncer/internal/logging"
	"testing"
)

func init() {
	logging.SetupForTests()
}

func newTestSyncAction(t *testing.T, action SyncActionType) SyncAction {
	tmpDir := t.TempDir()
	src := "test.file"
	fsutils.CreateDummyFile(filepath.Join(tmpDir, src), 10)
	return SyncAction{src: tmpDir, dst: tmpDir, filename: src, action: action}
}

func TestSyncActions(t *testing.T) {
	t.Run(string(COPY), func(t *testing.T) {
		action := newTestSyncAction(t, COPY)
		action.Run()
		_, err := os.Stat(filepath.Join(action.dst, action.filename))
		if err != nil {
			t.Fatalf("%s error %s", COPY, err)
		}
	})
	t.Run(string(MODIFY), func(t *testing.T) {
		action := newTestSyncAction(t, MODIFY)
		action.Run()
		_, err := os.Stat(filepath.Join(action.dst, action.filename))
		if err != nil {
			t.Fatalf("%s error %s", MODIFY, err)
		}
	})
	t.Run(string(REMOVE), func(t *testing.T) {
		action := newTestSyncAction(t, REMOVE)
		action.Run()
		_, err := os.Stat(filepath.Join(action.dst, action.filename))
		if err == nil {
			t.Fatalf("%s error %s", REMOVE, err)
		}
	})
}
