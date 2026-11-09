package middlewares

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	return c.ClientIP()
}

func ConfigureTrustedProxies(router *gin.Engine, raw string) error {
	if err := router.SetTrustedProxies(parseTrustedProxies(raw)); err != nil {
		return fmt.Errorf("invalid TRUSTED_PROXIES: %w", err)
	}

	return nil
}

func parseTrustedProxies(raw string) []string {
	var proxies []string

	for _, value := range strings.Split(raw, ",") {
		proxy := strings.TrimSpace(value)
		if proxy == "" {
			continue
		}
		proxies = append(proxies, proxy)
	}

	return proxies
}

func IPLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		clientIP := GetClientIP(c)

		c.Next()

		log.Printf("[IP] %s - %s %s - %d - %v",
			clientIP,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start))
	}
}
