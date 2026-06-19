// Package server konfiguriert und startet den HTTP-Server.
package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/myorg/myservice/config"
	"github.com/myorg/myservice/internal/handler"
	"github.com/myorg/myservice/internal/repository"
	"github.com/myorg/myservice/internal/service"
)

// New creates and configures an http.Server with a Gin router.
func New(cfg *config.Config, pool *pgxpool.Pool) *http.Server {
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())   // request logging
	router.Use(gin.Recovery()) // recover from panics

	// Dependency injection
	repo := repository.New(pool)
	svc := service.New(repo)
	h := handler.New(svc)

	// Frontend
	router.StaticFile("/", "./web/index.html")

	// Register routes
	registerRoutes(router, h)

	return &http.Server{
		Addr:    cfg.Addr(),
		Handler: router,
	}
}

// registerRoutes attaches all route groups to the router.
func registerRoutes(r *gin.Engine, h *handler.AutohausHandler) {
	// Health check
	r.GET("/health", healthHandler)

	rest := r.Group("/rest")
	{
		rest.GET("/:id", h.GetByID)
		rest.POST("", h.Create)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
