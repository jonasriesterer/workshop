// Package server konfiguriert und startet den HTTP-Server.
package server

import (
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/myorg/myservice/config"
	"github.com/myorg/myservice/internal/handler"
	"github.com/myorg/myservice/internal/middleware"
	"github.com/myorg/myservice/internal/repository"
	"github.com/myorg/myservice/internal/service"
)

// New creates and configures an http.Server with a Gin router.
func New(cfg *config.Config, pool *pgxpool.Pool, verifier *oidc.IDTokenVerifier) *http.Server {
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

	root := projectRoot()
	devH := handler.NewDev(pool,
		filepath.Join(root, "db", "migrations"),
		filepath.Join(root, "db", "seeds"),
	)

	// Frontend
	router.StaticFile("/", "./web/index.html")

	// Register routes
	registerRoutes(router, h, devH, cfg.KCClientID, verifier)

	return &http.Server{
		Addr:    cfg.Addr(),
		Handler: router,
	}
}

// registerRoutes attaches all route groups to the router.
func registerRoutes(r *gin.Engine, h *handler.AutohausHandler, devH *handler.DevHandler, clientID string, verifier *oidc.IDTokenVerifier) {
	// Health check
	r.GET("/health", healthHandler)

	jwtAuth := middleware.JWTAuth(verifier)

	rest := r.Group("/rest")
	{
		// Public read
		rest.GET("/:id", h.GetByID)

		// Requires admin or user role
		rest.POST("", jwtAuth, middleware.RequireRoles(clientID, "admin", "user"), h.Create)
		rest.PUT("/:id", jwtAuth, middleware.RequireRoles(clientID, "admin", "user"), h.Update)

		// Requires admin role only
		rest.DELETE("/:id", jwtAuth, middleware.RequireRoles(clientID, "admin"), h.Delete)
	}

	dev := r.Group("/dev")
	{
		dev.POST("/db_populate", devH.DbPopulate)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// projectRoot returns the absolute path to the project root,
// resolved relative to this source file at compile time.
func projectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	// server.go is at internal/server/ → two levels up is the project root
	return filepath.Join(filepath.Dir(filename), "..", "..")
}
