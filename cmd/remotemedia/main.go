package main

import (
	"github.com/scriptodude/remote-media/internal/remotemedia"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	remotemedia.StartServer()
}
