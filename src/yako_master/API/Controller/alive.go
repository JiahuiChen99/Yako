package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// IsAlive Handles YakoMaster aliveness
func IsAlive(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "alive"})
}
