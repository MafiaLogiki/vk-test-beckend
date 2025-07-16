// @title           VK Test Backend API
// @version         1.0
// @description     This is a backend service for VK test task.

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package main

import (
	"fmt"
	"net/http"
	"time"

	"marketplace-service/internal/auth"
	"marketplace-service/internal/config"
	"marketplace-service/internal/database"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/register"
	"marketplace-service/internal/token"

	docs "marketplace-service/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	secret = "vk-test"
)

func main() {
	l := logger.GetLogger()
	cfg := config.GetConfig(l)
	err := database.ConnectToDatabase(cfg, l)

	if err != nil {
		l.Fatal(err)
	}

	mux := http.NewServeMux()

	token := token.NewService(secret, time.Hour, l)
	
	docs.SwaggerInfo.BasePath = "/"
	mux.Handle("/api/v1/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))
	
	regHandler := register.NewHandler(l, token)
	regHandler.RegisterRoutes(mux)

	authHandler := auth.NewHandler(l, token)
	authHandler.RegisterService(mux)

	server := &http.Server {
		Addr: fmt.Sprintf("%s:%d", cfg.Listen.BindIp, cfg.Listen.Port),
		Handler: mux,
	}

	l.Info("Server is listening on port:", cfg.Listen.Port)
	if err := server.ListenAndServe(); err != nil {
		l.Fatal("Server failed to start:", err)
	}
}
