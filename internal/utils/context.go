package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUserID(c *gin.Context) (uint, error) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("не удалось извлечь userID из контекста")
	}
	userIDStr := fmt.Sprintf("%v", userIDRaw)
	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать userID: %v", err)
	}
	return uint(userIDUint64), nil
}
