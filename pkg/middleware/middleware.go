package middleware

import (
	"log"
	"net/http"
	"tender_management/controllers"
	"tender_management/pkg/utils"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func AutoMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("[ERROR] Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			if err.Error() == "token is expired" {
				log.Printf("[ERROR] Token expired: %v\n", err)
				controllers.HandleResponse(c, http.StatusUnauthorized, "Token expired")
			} else {
				log.Printf("[ERROR] Invalid token: %v\n", err)
				controllers.HandleResponse(c, http.StatusUnauthorized, "Invalid token")
			}
			c.Abort()
			return
		}
		alloved, err := e.Enforce(claims.Role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			controllers.HandleResponse(c, http.StatusInternalServerError, "Access denied")
			log.Println("[ERROR] Casbin enforcement error: ", err)
			c.Abort()
			return
		}

		if !alloved {
			controllers.HandleResponse(c,http.StatusForbidden, "Access denied")
			log.Printf("[INFO] Access denied for user: %v, Path: %s, Method: %s\n", claims.Role, c.Request.URL.Path, c.Request.Method)
			c.Abort()
			return
		}
		c.Next()
	}
}
