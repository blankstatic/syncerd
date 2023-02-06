package logging

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

const DefaultDisabledLogFile = ""

var Log *log.Logger

type Level struct {
	log.Level
}

func (l *Level) Set(val string) error {
	level, err := log.ParseLevel(val)
	if err != nil {
		return err
	}
	*l = Level{level}
	return nil
}

func Setup(logFormat string, logLevel Level, logFile string) {
	Log = log.New()

	Log.SetLevel(logLevel.Level)

	switch logFormat {
	case "text":
		Log.SetFormatter(&log.TextFormatter{})
	case "json":
		Log.SetFormatter(&log.JSONFormatter{
			PrettyPrint: false,
		})
	default:
	}

	Log.SetReportCaller(false)

	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			Log.Error("open file error: ", err)
		}
		w := io.MultiWriter(os.Stdout, f)
		Log.SetOutput(w)
	}
}

// tests used
func SetupForTests() {
	Setup("", Level{Level: log.ErrorLevel}, "")
}

func GetDefaultLogLevel() Level {
	return Level{Level: log.DebugLevel}
}
