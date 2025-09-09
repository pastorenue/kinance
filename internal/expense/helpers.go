package expense

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/pkg/middleware"
)

// getUserAndExpenseID extracts userID and expenseID from context and returns them, or writes an error response and returns false.
func getUserAndExpenseID(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	userIDVal, _ := c.Get(middleware.UserIDKey)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(500, common.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return uuid.Nil, uuid.Nil, false
	}
	expenseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, common.APIResponse{
			Success: false,
			Error:   "Invalid expense ID",
		})
		return uuid.Nil, uuid.Nil, false
	}
	return userID, expenseID, true
}

