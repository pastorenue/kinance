package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pastorenue/kinance/internal/common"
)

const UserIDKey = "user_id"

// AuthRequired now expects a token validation function to break the import cycle
func AuthRequired(validateToken func(string) (interface{}, error)) gin.HandlerFunc {
   return func(c *gin.Context) {
	   authHeader := c.GetHeader("Authorization")
	   if authHeader == "" {
		   c.JSON(http.StatusUnauthorized, common.APIResponse{
			   Success: false,
			   Error:   "Authorization header required",
		   })
		   c.Abort()
		   return
	   }

	   tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	   userID, err := validateToken(tokenString)
	   if err != nil {
		   c.JSON(http.StatusUnauthorized, common.APIResponse{
			   Success: false,
			   Error:   "Invalid token",
		   })
		   c.Abort()
		   return
	   }

	   c.Set(UserIDKey, userID)
	   c.Next()
   }
}
