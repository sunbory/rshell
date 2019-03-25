package rlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger

	logfileprefix = ".rshell/logs/rshell-"
	logfilesuffix = ".log"
)

func init() {
	rotate()

	now := fmt.Sprintf("%s", time.Now().UTC())
	now = strings.Replace(now, " ", "_", -1)
	now = strings.Replace(now, ":", "-", -1)

	logfile := logfileprefix + now + logfilesuffix
	lf, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("open the log file [%s] error: %v\n", logfile, err)
	}

	Debug = log.New(lf, "D: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(lf, "I: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(lf, "W: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, lf), "E: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func rotate() {
	if err := os.Mkdir(".rshell/logs", os.ModeDir); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("mkdir .rshell/logs error: %v\n", err)
		}
	}

	lfs, err := filepath.Glob(logfileprefix + "*")
	if err != nil {
		log.Fatalf("glob the log file [%s] error: %v\n", logfileprefix, err)
	}
	if len(lfs) > 10 {
		sort.Sort(sort.Reverse(sort.StringSlice(lfs)))
		for _, value := range lfs[10:] {
			os.Remove(value)
		}
	}
}
