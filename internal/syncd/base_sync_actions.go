package syncd

import (
	"os"
	"path/filepath"
	"syncer/internal/fsutils"
	"syncer/internal/logging"

	"github.com/sirupsen/logrus"
)

type SyncActionType string

const (
	COPY   SyncActionType = "copy"
	REMOVE SyncActionType = "remove"
	MODIFY SyncActionType = "modify"
)

type SyncAction struct {
	src      string
	dst      string
	filename string
	action   SyncActionType
}

func (action *SyncAction) Run() {
	var actionFn func() error

	switch action.action {
	case COPY:
		actionFn = action.copy
	case MODIFY:
		actionFn = action.modify
	case REMOVE:
		actionFn = action.remove
	default:
		logging.Log.Errorf("sync action '%v' not supported", action.action)
		return
	}
	if err := actionFn(); err != nil {
		logging.Log.Errorf("action error: %v", err)
	}
}

func (action *SyncAction) copy() error {
	fullSrc := filepath.Join(action.src, action.filename)
	fullDst := filepath.Join(action.dst, action.filename)
	logging.Log.WithFields(
		logrus.Fields{
			"size":   fsutils.GetFileSize(fullSrc),
			"action": action.action,
		},
	).Infof("copy from %v to %v", fullSrc, fullDst)

	if err := fsutils.MakeDirsForFile(fullDst); err != nil {
		return err
	}
	if isDir := fsutils.CheckDir(fullSrc); isDir {
		return nil
	}
	return fsutils.CopyFileContents(fullSrc, fullDst)
}

func (action *SyncAction) modify() error {
	return action.copy()
}

func (action *SyncAction) remove() error {
	fullDst := filepath.Join(action.dst, action.filename)
	logging.Log.WithFields(
		logrus.Fields{
			"size":   fsutils.GetFileSize(fullDst),
			"action": action.action,
		},
	).Infof("remove %v", fullDst)

	return os.Remove(fullDst)
}
