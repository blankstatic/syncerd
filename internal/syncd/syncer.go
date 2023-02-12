package syncd

type Syncer interface {
	Cleanup()
	Startup()
	RunSyncLoop()
}

func SyncLoop(s Syncer) {
	defer s.Cleanup()
	s.Startup()
	s.RunSyncLoop()
}
