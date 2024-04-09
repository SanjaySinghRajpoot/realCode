package controller

import (
	"net/http"

	"github.com/SanjaySinghRajpoot/realCode/utils"
	"github.com/gin-gonic/gin"
)

func CompileHandler(c *gin.Context) {

	// Code should be formatted in the JSON format
	var request struct {
		Language string `json:"language"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check what is the language of the code and call the appropitrate service
	var outputResult string
	if request.Language == "python" {
		output, err := utils.CompileCodePython(request.Code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		outputResult = output
	} else if request.Language == "golang" {
		output, err := utils.CompileCodeGo(request.Code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		outputResult = output
	}

	c.JSON(http.StatusOK, gin.H{"output": outputResult})
}
