package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRoles ensures the authenticated user holds at least one of the given
// client roles.  Roles are read from the Keycloak claim:
//
//	resource_access.<clientID>.roles
//
// Must be used after JWTAuth which populates the claims in the context.
func RequireRoles(clientID string, roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		raw, exists := c.Get(claimsKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no claims in context"})
			return
		}

		claims, ok := raw.(map[string]any)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid claims type"})
			return
		}

		userRoles := extractClientRoles(claims, clientID)

		for _, r := range userRoles {
			if _, ok := allowed[r]; ok {
				c.Next()
				return
			}
		}

		slog.Warn("access denied – missing role", "required", roles, "clientID", clientID)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
	}
}

// extractClientRoles reads resource_access[clientID].roles from the JWT claims.
func extractClientRoles(claims map[string]any, clientID string) []string {
	resourceAccess, ok := claims["resource_access"].(map[string]any)
	if !ok {
		return nil
	}

	clientAccess, ok := resourceAccess[clientID].(map[string]any)
	if !ok {
		return nil
	}

	rawRoles, ok := clientAccess["roles"].([]any)
	if !ok {
		return nil
	}

	roles := make([]string, 0, len(rawRoles))
	for _, r := range rawRoles {
		if s, ok := r.(string); ok {
			roles = append(roles, s)
		}
	}
	return roles
}
