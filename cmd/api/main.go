package main

import (
	"main/core/config"
	"main/core/database/sqlite"
	"main/core/logger"
	"main/core/middleware"
	"main/core/server"
	"os"
	"os/signal"
)

func init() {
	config.Init()
	logger.Init()
	sqlite.Init()
	middleware.Init()
}

func main() {
	go server.Run()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	server.Shutdown()
	logger.Get().Info("Server stopped")
}
