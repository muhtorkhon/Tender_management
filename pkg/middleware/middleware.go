package middleware

import (
	"log"
	"net/http"
	"tender_management/controllers"
	"tender_management/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AutoMiddleware(role string) gin.HandlerFunc {
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

		if claims.Role != role {
			log.Printf("[ERROR] Access denied for role: %s, required: %s\n", claims.Role, role)
			controllers.HandleResponse(c, http.StatusForbidden, "Access denied")
			c.Abort()
			return
		}

		log.Printf("[INFO] Authorized request: Email: %s, Role: %s\n", claims.Email, claims.Role)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}
