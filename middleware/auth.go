package middleware

import (
	"net/http"
	"strings"

	"survey-backend/internal/user"
	"survey-backend/pkg/jwt"
	"survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	r := &response.Response[map[string]string]{}
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, r.Fail("missing token"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, r.Fail("invalid token format"))
			c.Abort()
			return
		}

		tokenStr := parts[1]
		email, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, r.Fail("invalid token"))
			c.Abort()
			return
		}

		var u user.User
		if err := user.FindUserByEmail(email, &u); err != nil {
			c.JSON(http.StatusUnauthorized, r.Fail("user not found"))
			c.Abort()
			return
		}

		c.Set("currentUser", u)
		c.Next()
	}
}
