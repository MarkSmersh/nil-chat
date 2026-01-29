package main

import (
	"log"
	"log/slog"

	"github.com/MarkSmersh/nil-chat/api"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	api.Init()
}
