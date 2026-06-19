// Package handler enthält die HTTP-Handler für alle Endpunkte.
package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/myorg/myservice/internal/model"
	"github.com/myorg/myservice/internal/repository"
	"github.com/myorg/myservice/internal/service"
)

const (
	msgInvalidID       = "invalid id: %s"
	msgNotFound        = "autohaus with id %d not found"
	msgInternalError   = "internal server error"
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

// GetByID handles GET /rest/:id with conditional GET support via ETag/If-None-Match.
func (h *AutohausHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(msgInvalidID, idStr)})
		return
	}

	ah, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf(msgNotFound, id)})
			return
		}
		slog.Error("GetByID failed", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msgInternalError})
		return
	}

	etag := fmt.Sprintf(`"%d"`, ah.Version)
	c.Header("ETag", etag)

	if c.GetHeader("If-None-Match") == etag {
		slog.Info("GetByID not modified", "id", id)
		c.Status(http.StatusNotModified)
		return
	}

	slog.Info("GetByID", "id", id)
	c.JSON(http.StatusOK, ah)
}

// Update handles PUT /rest/:id — replaces the resource; requires If-Match for optimistic locking.
func (h *AutohausHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(msgInvalidID, idStr)})
		return
	}

	ifMatch := c.GetHeader("If-Match")
	if ifMatch == "" {
		c.JSON(http.StatusPreconditionRequired, gin.H{"error": "If-Match header required"})
		return
	}
	versionStr := strings.Trim(ifMatch, `"`)
	version, err := strconv.ParseUint(versionStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid If-Match header"})
		return
	}

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

	input.ID = uint(id)
	input.Version = uint(version)

	if err := h.svc.Update(c.Request.Context(), &input); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf(msgNotFound, id)})
		case errors.Is(err, repository.ErrConflict):
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": "version conflict"})
		default:
			slog.Error("Update failed", "id", id, "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msgInternalError})
		}
		return
	}

	slog.Info("Update", "id", id)
	c.Status(http.StatusNoContent)
}

// Delete handles DELETE /rest/:id
func (h *AutohausHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf(msgInvalidID, idStr)})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf(msgNotFound, id)})
			return
		}
		slog.Error("Delete failed", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msgInternalError})
		return
	}

	slog.Info("Delete", "id", id)
	c.Status(http.StatusNoContent)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": msgInternalError})
		return
	}

	location := fmt.Sprintf("/rest/%d", id)
	slog.Info("Create", "id", id, "location", location)
	c.Header("Location", location)
	c.Status(http.StatusCreated)
}
