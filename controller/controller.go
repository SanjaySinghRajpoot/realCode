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
		UserID   uint   `json:"user_id"`
		// add user id here
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check what is the language of the code and call the appropitrate service
	var outputResult string
	var err error

	if request.Language == "python" {

		outputResult, err = utils.CompileCodePython(request.Code, request.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if request.Language == "golang" {

		outputResult, err = utils.CompileCodeGo(request.Code, request.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// here we will add the code to comile the code

	// if request.Language == "python" {

	// 	msg, err := kafka.Producer("python", request, kafka.KafkaProducer)
	// 	if err != nil {
	// 		fmt.Println(msg)
	// 		formatError.InternalServerError(c, err)
	// 		return
	// 	}

	// 	outputResult = msg
	// } else if request.Language == "golang" {

	// 	msg, err := kafka.Producer("golang", request, kafka.KafkaProducer)
	// 	if err != nil {
	// 		fmt.Println(msg)
	// 		formatError.InternalServerError(c, err)
	// 		return
	// 	}

	// 	output, err := utils.CompileCodeGo(request.Code)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	outputResult = output
	// }

	c.JSON(http.StatusOK, gin.H{"output": outputResult})
}
