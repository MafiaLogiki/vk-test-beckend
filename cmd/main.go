package main

import (
	"marketplace-service/internal/database"
	"marketplace-service/internal/config"
)

func main() {
	database.InitDatabase()
}
