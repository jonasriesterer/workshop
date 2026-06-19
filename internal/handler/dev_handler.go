package handler

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/myorg/myservice/internal/seed"
)

// DevHandler provides development-only endpoints.
type DevHandler struct {
	pool         *pgxpool.Pool
	migrationDir string
	csvDir       string
}

// NewDev creates a DevHandler with the given pool and filesystem paths.
func NewDev(pool *pgxpool.Pool, migrationDir, csvDir string) *DevHandler {
	return &DevHandler{pool: pool, migrationDir: migrationDir, csvDir: csvDir}
}

// DbPopulate handles POST /dev/db_populate — drops, recreates and reseeds the database.
func (h *DevHandler) DbPopulate(c *gin.Context) {
	ctx := c.Request.Context()

	for _, file := range []string{"001_drop.sql", "002_create.sql"} {
		path := filepath.Join(h.migrationDir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			slog.Error("DbPopulate: read migration", "file", file, "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "migration file not found"})
			return
		}
		if _, err := h.pool.Exec(ctx, string(content)); err != nil {
			slog.Error("DbPopulate: exec migration", "file", file, "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "migration failed"})
			return
		}
	}

	if err := seed.Load(ctx, h.pool, h.csvDir); err != nil {
		slog.Error("DbPopulate: seed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "seeding failed"})
		return
	}

	slog.Info("DbPopulate: database reset and seeded")
	c.Status(http.StatusOK)
}
