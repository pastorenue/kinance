package middleware

import (
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs HTTP requests
// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// Logger logs HTTP requests with colors
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()
		
		duration := time.Since(start)
		status := c.Writer.Status()
		
		// Color based on status code
		var statusColor string
		switch {
		case status >= 500:
			statusColor = colorRed
		case status >= 400:
			statusColor = colorYellow
		case status >= 300:
			statusColor = colorCyan
		case status >= 200:
			statusColor = colorGreen
		default:
			statusColor = colorReset
		}
		
		// Color based on method
		var methodColor string
		switch c.Request.Method {
		case "GET":
			methodColor = colorBlue
		case "POST":
			methodColor = colorGreen
		case "PUT":
			methodColor = colorYellow
		case "DELETE":
			methodColor = colorRed
		default:
			methodColor = colorCyan
		}
		
		log.Printf("%s%s%s %s %s%d%s %v",
			methodColor, c.Request.Method, colorCyan,
			c.Request.URL.Path,
			statusColor, status, colorReset,
			duration)
	}
}

// Recovery recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v\n%s", err, debug.Stack())
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS handles cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		
		c.Next()
	}
}

// RateLimit limits requests per IP
func RateLimit(rateLimit int) gin.HandlerFunc {
	type client struct {
		count      int
		lastAccess time.Time
	}
	
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
		limit   = rateLimit // requests per minute
	)
	
	// Cleanup old entries every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastAccess) > time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		mu.Lock()
		cl, exists := clients[ip]
		if !exists {
			cl = &client{}
			clients[ip] = cl
		}
		
		// Reset counter if minute has passed
		if time.Since(cl.lastAccess) > time.Minute {
			cl.count = 0
		}
		
		cl.count++
		cl.lastAccess = time.Now()
		
		if cl.count > limit {
			mu.Unlock()
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		mu.Unlock()
		
		c.Next()
	}
}