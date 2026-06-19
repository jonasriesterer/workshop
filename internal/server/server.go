package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myorg/myservice/config"
)

// New creates and configures an http.Server with a Gin router.
func New(cfg *config.Config) *http.Server {
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())   // request logging
	router.Use(gin.Recovery()) // recover from panics

	// Register routes
	registerRoutes(router)

	return &http.Server{
		Addr:    cfg.Addr(),
		Handler: router,
	}
}

// registerRoutes attaches all route groups to the router.
// Add your API routes here as the service grows.
func registerRoutes(r *gin.Engine) {
	// Health check — used by load balancers, Docker, k8s liveness probes
	r.GET("/health", healthHandler)

	// Future API groups go here, e.g.:
	// v1 := r.Group("/api/v1")
	// {
	//     v1.GET("/users", ...)
	// }
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
