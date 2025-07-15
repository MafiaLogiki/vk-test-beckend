package main

import (
	"marketplace-service/internal/config"
	"marketplace-service/internal/database"
	"marketplace-service/internal/logger"
)

func main() {
	l := logger.GetLogger()
	cfg := config.GetConfig(l)
	err := database.ConnectToDatabase(cfg, l)

	if err != nil {
		l.Fatal(err)
	}
}
