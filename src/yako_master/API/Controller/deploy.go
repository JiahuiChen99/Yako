package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// UploadApp handles the file that the user wants
// to deploy in the cluster
func UploadApp(c *gin.Context) {

	// File uploaded and stored
	c.JSON(http.StatusOK, map[string]string{"status": "uploaded successfully"})
}
