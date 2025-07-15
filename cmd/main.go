package main

import (
	"marketplace-service/internal/config"
	"marketplace-service/internal/database"
	"marketplace-service/internal/logger"
)

func main() {
	l := logger.GetLogger()
	cfg := config.GetConfig(l)
	database.ConnectToDatabase(cfg)
}
