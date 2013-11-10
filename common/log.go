package common

import (
	"log"
	"os"
)

var LogErr, LogWarn, LogInfo *log.Logger

type logWriter struct {
}

func (lw *logWriter) Write(p []byte) {

}

func init() {
	LogErr = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Llongfile)
	LogWarn = log.New(os.Stderr, "WARN: ", log.LstdFlags|log.Llongfile)
	LogInfo = log.New(os.Stderr, "INFO: ", log.LstdFlags|log.Llongfile)
}
