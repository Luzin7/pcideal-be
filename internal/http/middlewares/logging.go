package middlewares

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		if isValidIP(ip) {
			return ip
		}
	}

	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" && isValidIP(realIP) {
		return realIP
	}

	return c.ClientIP()
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IPLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		clientIP := GetClientIP(c)

		c.Next()

		log.Printf("%s - %s %s - %d - %v",
			clientIP,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start))
	}
}
