package main

import (
	"main/core/config"
	"main/core/database/sqlite"
	"main/core/logger"
	"main/core/server"
)

func init() {
	config.Init()
	logger.Init()
	sqlite.Init()
}

func main() {
	server.Run()
}
