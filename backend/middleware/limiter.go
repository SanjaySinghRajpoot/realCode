package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/SanjaySinghRajpoot/realCode/backend/utils/localredis"
	"github.com/gin-gonic/gin"
)

const (
	maxRequests     = 25
	perMinutePeriod = 1 * time.Minute
)

var (
	mutex = &sync.Mutex{}
)

func RateLimiter(context *gin.Context) {
	ip := context.ClientIP()
	mutex.Lock()
	defer mutex.Unlock()

	isPresent, _ := localredis.CheckIPAddressKey(ip)

	fmt.Println(isPresent)

	count := 0

	if isPresent {
		// get from redis cache

		getCount, err := localredis.GetIPAddress(ip)

		if err != nil {
			fmt.Printf("Failed to Get the Redis Cache: %d", getCount)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		fmt.Println(getCount)

		if getCount >= maxRequests {
			context.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		count = getCount
	}

	count = count + 1

	fmt.Println(ip)

	msg, err := localredis.SetIPAddress(ip, count)
	if err != nil {
		fmt.Printf("Failed to Set the Redis Cache: %s", msg)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	time.AfterFunc(perMinutePeriod, func() {
		mutex.Lock()
		defer mutex.Unlock()

		count, err := localredis.GetIPAddress(ip)

		if err != nil {
			fmt.Printf("Failed to Get the Redis Cache CRON: %d", count)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		count = count - 1

		msg, err := localredis.SetIPAddress(ip, count)
		if err != nil {
			fmt.Printf("Failed to Set the Redis Cache CRON: %s", msg)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

	})

	context.Next()
}