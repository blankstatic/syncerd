package syncd

import (
	"context"
	"sync"
	"syncer/pkg/fsutils"
	"syncer/pkg/lck"
	"syncer/pkg/logging"
	"time"
)

const FullSyncInterval = 30 * time.Second

type BaseSync struct {
	sync.Mutex
	ctx           context.Context
	src           string
	dst           string
	isInitialized bool
	syncLocked    bool
	srcData       fsutils.MD5Files
	dstData       fsutils.MD5Files
	watcher       *SyncWatcher
	force         bool
	interval      time.Duration
}

func NewBaseSync(ctx context.Context, src, dst string, force bool, interval time.Duration) Syncer {
	return &BaseSync{ctx: ctx, src: src, dst: dst, force: force, interval: interval}
}

func (bs *BaseSync) RunSyncLoop() {
	fullSyncTicker := time.NewTicker(bs.interval)
	defer fullSyncTicker.Stop()

	logging.Log.Info("sync start")
	bs.fillSyncMap()
	bs.fullSync()

	changesCh := make(chan *SyncAction)
	bs.watcher = &SyncWatcher{
		ch:  changesCh,
		ctx: bs.ctx,
		src: bs.src,
		dst: bs.dst,
	}
	go bs.watcher.Watch()

syncLoop:
	for {
		select {
		case <-bs.ctx.Done():
			logging.Log.Info("sync stop")
			break syncLoop
		case <-fullSyncTicker.C:
			bs.watcher.Disable()
			logging.Log.Infof("[%v] sync...", bs.interval)
			bs.fillSyncMap()
			bs.fullSync()
			bs.watcher.Enable()
		case change := <-changesCh:
			bs.execSyncActions(*change)
		}
	}
	logging.Log.Info("sync done")
}

func (bs *BaseSync) Cleanup() {
	if !bs.force {
		lck.Unlock()
		logging.Log.Info("sync cleanup")
	}
}

func (bs *BaseSync) Startup() {
	if !bs.force {
		isLocked, lockFile := lck.IsLocked()
		if isLocked {
			logging.Log.Fatalf("app is locked by %s", lockFile)
		}
		logging.Log.Info("sync startup")
	}
}

func (bs *BaseSync) setData(srcData *fsutils.MD5Files, dstData *fsutils.MD5Files) {
	bs.Lock()
	defer bs.Unlock()

	bs.srcData = *srcData
	bs.dstData = *dstData

	if !bs.isInitialized {
		bs.isInitialized = true
		logging.Log.Info("src total: ", len(bs.srcData))
		logging.Log.Info("dst total: ", len(bs.dstData))
	}
}

func (bs *BaseSync) fillSyncMap() {
	srcData, err := fsutils.MD5All(bs.src, true)
	if err != nil {
		logging.Log.Error(err)
	}
	dstData, err := fsutils.MD5All(bs.dst, true)
	if err != nil {
		logging.Log.Error(err)
	}
	bs.setData(&srcData, &dstData)
}

func (bs *BaseSync) fullSync() {
	bs.Lock()
	defer bs.Unlock()

	if bs.syncLocked {
		logging.Log.Warning("sync locked")
		return
	}

	bs.syncLocked = true
	defer func() { bs.syncLocked = false }()

	actions := bs.getSyncActions()
	bs.execSyncActions(actions...)
}

func (bs *BaseSync) getSyncActions() []SyncAction {
	actions := []SyncAction{}

	for srcFilename, srcHash := range bs.srcData {
		if _, exist := bs.dstData[srcFilename]; !exist {
			actions = append(actions, SyncAction{
				src:      bs.src,
				dst:      bs.dst,
				filename: srcFilename,
				action:   COPY,
			})
			bs.dstData[srcFilename] = srcHash
		} else {
			if srcHash != bs.dstData[srcFilename] {
				actions = append(actions, SyncAction{
					src:      bs.src,
					dst:      bs.dst,
					filename: srcFilename,
					action:   MODIFY,
				})
				bs.dstData[srcFilename] = srcHash
			}
		}
	}

	for dstFilename := range bs.dstData {
		if _, exist := bs.srcData[dstFilename]; !exist {
			actions = append(actions, SyncAction{
				dst:      bs.dst,
				filename: dstFilename,
				action:   REMOVE,
			})
			delete(bs.dstData, dstFilename)
		}
	}

	return actions
}

func (bs *BaseSync) execSyncActions(actions ...SyncAction) {
	if len(actions) == 0 {
		return
	}

	logging.Log.Infof("sync %v changes", len(actions))
	for _, action := range actions {
		action.Run()
	}
}
