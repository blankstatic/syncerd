package syncd

import (
	"fmt"
	"path/filepath"
	"syncer/pkg/fsutils"
	"testing"
)

func BenchmarkCopy(b *testing.B) {
	src := b.TempDir()
	dst := filepath.Join(src, "dst")
	prepDst := filepath.Join(src, "some.file")
	err := fsutils.MakeDirsForFile(prepDst)
	if err != nil {
		b.Fatalf("error create dir %s", dst)
	}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		filename := fmt.Sprint(i)
		action := &SyncAction{src: src, dst: dst, filename: filename, action: COPY}
		fsutils.CreateDummyFile(filepath.Join(src, filename), 1)

		b.StartTimer()
		action.Run()
	}
}
