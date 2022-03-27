package Controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yako/src/model"
	"yako/src/utils/directory_util"
	"yako/src/yako_master/API/utils"
)

// UploadApp handles the file that the user wants
// to deploy in the cluster
func UploadApp(c *gin.Context) {
	file, formErr := c.FormFile("app")
	if formErr != nil {
		err := utils.BadRequestError(formErr.Error())
		c.JSON(err.Status, err)
		return
	}

	// Check if YakoMaster's working directory is available
	directory_util.WorkDir("yakomaster")

	// Save the file on the server
	if saveErr := c.SaveUploadedFile(file, "/usr/yakomaster/"+file.Filename); saveErr != nil {
		err := utils.InternalServerError(saveErr.Error())
		c.JSON(err.Status, err)
		return
	}

	// Get the app's resources configuration
	var config model.Config
	if err := c.ShouldBind(&config); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	// File uploaded and stored
	c.JSON(http.StatusOK, map[string]string{"status": "uploaded successfully"})
}
