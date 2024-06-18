package controller

import (
	"fmt"
	"net/http"

	"github.com/SanjaySinghRajpoot/realCode/utils"
	"github.com/SanjaySinghRajpoot/realCode/utils/formatError"
	"github.com/SanjaySinghRajpoot/realCode/utils/kafka"
	"github.com/gin-gonic/gin"
)

func CompileHandler(c *gin.Context) {

	// Code should be formatted in the JSON format
	var request struct {
		Language string `json:"language"`
		Code     string `json:"code"`
		UserID   uint   `json:"user_id"`
		// add user id here
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check what is the language of the code and call the appropitrate service
	var outputResult string
	if request.Language == "python" {

		msg, err := kafka.Producer("python", request, kafka.KafkaProducer)
		if err != nil {
			fmt.Println(msg)
			formatError.InternalServerError(c, err)
			return
		}

		outputResult = msg
	} else if request.Language == "golang" {

		msg, err := kafka.Producer("golang", request, kafka.KafkaProducer)
		if err != nil {
			fmt.Println(msg)
			formatError.InternalServerError(c, err)
			return
		}

		output, err := utils.CompileCodeGo(request.Code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		outputResult = output
	}

	c.JSON(http.StatusOK, gin.H{"output": outputResult})
}
