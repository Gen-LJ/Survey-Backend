package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware answers browser preflights and echoes back an allowed origin.
// A Flutter mobile build does not need this; a Flutter web build does, because
// the browser will refuse any cross-origin call without it.
//
// Pass the origins from CORS_ALLOWED_ORIGINS. An empty list disables CORS
// entirely, which is the right default for a mobile-only API. The single entry
// "*" allows any origin, but note that browsers reject a wildcard on
// credentialed requests, so a real origin list is better once you know it.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowAll := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"

	allowed := make(map[string]bool, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[strings.ToLower(strings.TrimRight(origin, "/"))] = true
	}

	maxAge := strconv.Itoa(int((12 * time.Hour).Seconds()))

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Not a cross-origin request; nothing to negotiate.
		if origin == "" {
			c.Next()
			return
		}

		switch {
		case allowAll:
			c.Header("Access-Control-Allow-Origin", "*")
		case allowed[strings.ToLower(strings.TrimRight(origin, "/"))]:
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			// The response varies per origin, so caches must not share it.
			c.Header("Vary", "Origin")
		default:
			// Unknown origin: send no CORS headers and let the browser block it.
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Max-Age", maxAge)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
