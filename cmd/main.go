// @title           VK Test Backend API
// @version         1.0
// @description     This is a backend service for VK test task.

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

package main

import (
	"fmt"
	"net/http"
	"time"

	"marketplace-service/internal/announcements"
	"marketplace-service/internal/auth"
	"marketplace-service/internal/config"
	"marketplace-service/internal/database"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/register"
	"marketplace-service/internal/store"
	"marketplace-service/internal/token"

	docs "marketplace-service/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)


func main() {
	l := logger.GetLogger()
	cfg := config.GetConfig(l)
	db, err := database.ConnectToDatabase(cfg, l)

	if err != nil {
		l.Fatal(err)
	}

	mux := http.NewServeMux()

	token := token.NewService(cfg.Secret, time.Hour, l)
	
	docs.SwaggerInfo.BasePath = "/"
	mux.Handle("/api/v1/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))
	

	userStore := store.NewPostgresUserStore(db)

	regHandler := register.NewHandler(userStore, l, token)
	regHandler.RegisterRoutes(mux)

	authHandler := auth.NewHandler(userStore, l, token)
	authHandler.RegisterService(mux)

	announcementStore := store.NewPostgresAnnouncementsStore(db)

	announcementsHandler := announcements.NewHandler(announcementStore, l, token)
	announcementsHandler.RegisterService(mux)

	server := &http.Server {
		Addr: fmt.Sprintf("%s:%d", cfg.Listen.BindIp, cfg.Listen.Port),
		Handler: mux,
	}

	l.Info("Server is listening on port:", cfg.Listen.Port)
	if err := server.ListenAndServe(); err != nil {
		l.Fatal("Server failed to start:", err)
	}
}
