package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/sahilrana27582/sentinel-ai/sre-api/internal/config"
	"github.com/sahilrana27582/sentinel-ai/sre-api/internal/handler"
	"github.com/sahilrana27582/sentinel-ai/sre-api/internal/routes"
)

func main() {
	cfg := config.Load()
	log.Println(cfg.String())

	r := gin.Default()

	// Initialize handlers
	h := handler.New(cfg)

	// Setup routes
	routes.Setup(r, h)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.APIHost, cfg.APIPort),
		Handler: r,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), config.RequiredStartupDuration)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
	log.Println("server stopped")
}
