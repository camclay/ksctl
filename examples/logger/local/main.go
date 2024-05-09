package main

import (
	"os"

	"github.com/ksctl/ksctl/pkg/logger"
	"github.com/ksctl/ksctl/pkg/resources"
)

func main() {
	f, _ := os.OpenFile("/tmp/demo.log", os.O_CREATE|os.O_RDWR, 0755)
	var logFile resources.LoggerFactory = logger.NewStructuredLogger(-1, f)
	logFile.Print("Example", "key", "value")
	logFile.Debug("HelloStdout")
	logFile.Box("hello", "world!")

	var log resources.LoggerFactory = logger.NewStructuredLogger(-1, os.Stdout)
	log.Print("Example", "key", "value")
	log.Debug("HelloStdout")
}
