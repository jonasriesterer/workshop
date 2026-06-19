// Package middleware enthält HTTP-Middleware für Authentifizierung und Autorisierung.
package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

const claimsKey = "oidc_claims"

// JWTAuth validates an OIDC Bearer token from the Authorization header.
// On success the raw claims map is stored under claimsKey in the Gin context.
func JWTAuth(verifier *oidc.IDTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed Authorization header"})
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, err := verifier.Verify(c.Request.Context(), rawToken)
		if err != nil {
			slog.Warn("token verification failed", "err", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		var claims map[string]any
		if err := idToken.Claims(&claims); err != nil {
			slog.Error("failed to parse token claims", "err", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse token claims"})
			return
		}

		c.Set(claimsKey, claims)
		c.Next()
	}
}
