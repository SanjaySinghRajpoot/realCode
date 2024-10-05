package main

import (
	"net/http"

	"github.com/SanjaySinghRajpoot/realCode/controller"
	"github.com/SanjaySinghRajpoot/realCode/middleware"
	"github.com/SanjaySinghRajpoot/realCode/utils/redis"
	"github.com/gin-gonic/gin"
)

func HomepageHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to Real Code."})
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

	// connect to DB
	// config.ConnectDB()

	password := "12345678"

	// Redis Cache Setup
	redis.RedisClient = redis.SetUpRedis(password)

	// var err error
	// kafka.KafkaProducer, err = kafka.InitializeProducer()

	// if err != nil {
	// 	fmt.Printf("Failed to create producer: %s\n", err.Error())
	// 	return
	// }

	router := gin.Default()

	router.Use(CORS())

	router.GET("/", HomepageHandler)

	router.POST("/compile", middleware.RateLimiter, controller.CompileHandler)
	router.Run(":8080")

}
