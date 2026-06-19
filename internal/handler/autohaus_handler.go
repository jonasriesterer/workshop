// Package handler enthält die HTTP-Handler für alle Endpunkte.
package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/myorg/myservice/internal/model"
	"github.com/myorg/myservice/internal/repository"
	"github.com/myorg/myservice/internal/service"
)

// AutohausHandler handles HTTP requests for the Autohaus resource.
type AutohausHandler struct {
	svc      service.AutohausService
	validate *validator.Validate
}

// New creates a new AutohausHandler.
func New(svc service.AutohausService) *AutohausHandler {
	return &AutohausHandler{
		svc:      svc,
		validate: validator.New(),
	}
}

// GetByID handles GET /rest/:id
func (h *AutohausHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid id: %s", idStr)})
		return
	}

	ah, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("autohaus with id %d not found", id)})
			return
		}
		slog.Error("GetByID failed", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	slog.Info("GetByID", "id", id)
	c.JSON(http.StatusOK, ah)
}

// Create handles POST /rest
func (h *AutohausHandler) Create(c *gin.Context) {
	var input model.Autohaus
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errs := make([]string, 0, len(ve))
			for _, e := range ve {
				errs = append(errs, fmt.Sprintf("field '%s' failed '%s'", e.Field(), e.Tag()))
			}
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": errs})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.svc.Create(c.Request.Context(), &input)
	if err != nil {
		slog.Error("Create failed", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	location := fmt.Sprintf("/rest/%d", id)
	slog.Info("Create", "id", id, "location", location)
	c.Header("Location", location)
	c.Status(http.StatusCreated)
}
