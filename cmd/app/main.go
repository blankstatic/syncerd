package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syncer/internal/fsutils"
	"syncer/internal/logging"
	"syncer/internal/syncd"
	"syscall"
	"time"
)

var TerminationSignals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

type Cfg struct {
	src, dst         string
	logFormat        string
	logLevel         logging.Level
	logFile          string
	force            bool
	fullSyncInterval time.Duration
}

func initCfg() *Cfg {
	cfg := Cfg{logLevel: logging.GetDefaultLogLevel()}
	flag.StringVar(&cfg.src, "src", "", "source sync dir")
	flag.StringVar(&cfg.dst, "dst", "", "destination sync dir")
	flag.StringVar(&cfg.logFile, "log", logging.DefaultDisabledLogFile, "log file")
	flag.StringVar(&cfg.logFormat, "format", "json", "output log format [text, json]")
	flag.Var(&cfg.logLevel, "level", "log level")
	flag.BoolVar(&cfg.force, "force", false, "force mode (ignore lockfile)")
	flag.DurationVar(&cfg.fullSyncInterval, "interval", syncd.FullSyncInterval, "full sync interval")

	flag.Usage = func() { fmt.Println("Please specify args.") }

	flag.Parse()

	if len(cfg.src) == 0 || len(cfg.dst) == 0 {
		flag.Usage()
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("--%v (default: '%v')\n", f.Name, f.Value)
			fmt.Fprintf(os.Stderr, "    %v\n", f.Usage)
		})
		os.Exit(1)
	}
	cfg.src = fsutils.ResolveAbsPath(cfg.src)
	cfg.dst = fsutils.ResolveAbsPath(cfg.dst)
	if cfg.src == cfg.dst {
		logging.Log.Fatal("src and dst should be different")
	}
	return &cfg
}

func run() {
	cfg := initCfg()

	logging.Setup(cfg.logFormat, cfg.logLevel, cfg.logFile)

	dirsErr := fsutils.CheckDirs(cfg.src, cfg.dst)
	if dirsErr != nil {
		logging.Log.WithError(dirsErr).Fatal("app exiting")
	}

	logging.Log.Info("app start")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalTermination := make(chan os.Signal, 1)
	signal.Notify(signalTermination, TerminationSignals...)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s := <-signalTermination
		cancel()
		logging.Log.WithError(errors.New(s.String())).Info("app terminating")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bs := syncd.NewBaseSync(ctx, cfg.src, cfg.dst, cfg.force, cfg.fullSyncInterval)
		syncd.SyncLoop(bs)
	}()

	wg.Wait()
}

func main() {
	run()
}
