package syncd

import (
	"context"
	"path/filepath"
	"reflect"
	"syncer/internal/fsutils"
	"syncer/internal/logging"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	logging.SetupForTests()
}

func newTestSyncWatcher(t *testing.T) *SyncWatcher {
	tmpDir := t.TempDir()
	ch := make(chan *SyncAction)
	return &SyncWatcher{src: tmpDir, dst: tmpDir, ch: ch}
}

func TestSyncWatcherToggle(t *testing.T) {
	w := newTestSyncWatcher(t)
	assert.False(t, w.isLocked())

	w.Disable()
	assert.True(t, w.isLocked())

	w.Enable()
	assert.False(t, w.isLocked())
}

func TestSyncWatcher(t *testing.T) {
	w := newTestSyncWatcher(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.ctx = ctx

	t.Run("format name", func(t *testing.T) {
		testFilename := "test_case.file"
		fmtFilename := w.formatName(w.src + testFilename)

		assert.Equal(t, fmtFilename, testFilename)
	})

	t.Run("subdirs", func(t *testing.T) {
		type test struct {
			src  string
			want []string
		}
		dirs := filepath.Join(w.src, "dir", "subdir1", "subdir2", "some.file")
		err := fsutils.MakeDirsForFile(dirs)
		if err != nil {
			t.Fatalf("error creating dir %s", dirs)
		}

		tests := []test{
			{src: w.src, want: []string{
				w.src,
				filepath.Join(w.src, "dir"),
				filepath.Join(w.src, "dir", "subdir1"),
				filepath.Join(w.src, "dir", "subdir1", "subdir2"),
			}},
		}

		for _, tc := range tests {
			got := w.getSubdirs()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		}
	})

	t.Run("send", func(t *testing.T) {
		action := &SyncAction{action: "unknown"}
		w.send(action)

		select {
		case <-w.ch:
			t.Fatal("send is blocked")
		default:
		}
	})

	t.Run("watch", func(t *testing.T) {
		func() {
			<-time.After(time.Millisecond * 10)
			cancel()
		}()
		w.Watch()
	})
}
