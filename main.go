package main

import (
	"bytes"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func compileCode(code string) (string, error) {
	cmd := exec.Command("bash", "-c", "echo '"+code+"' | python3 -c 'import sys; exec(sys.stdin.read())'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func compileHandler(c *gin.Context) {
	var request struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := compileCode(request.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": output})
}

func main() {
	router := gin.Default()
	router.POST("/compile", compileHandler)
	router.Run(":8080")

}
