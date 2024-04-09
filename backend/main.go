package main

import (
	"net/http"

	"github.com/SanjaySinghRajpoot/realCode/controller"
	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real code"})
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set Headers
		// c.Writer.Header().Set("Access-Control-Allow-Headers:", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()

	router.Use(CORS())

	router.GET("/", HomepageHandler)

	router.POST("/compile", controller.CompileHandler)
	router.Run(":8080")

}
