package middleware

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// Logger mencatat request masuk & keluar.
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()
        method := c.Request.Method
        path := c.Request.URL.Path

        log.Printf("[%d] %s %s (%s)", status, method, path, latency)
    }
}

// CORS sederhana dan aman—bisa dibuat configurable nanti.
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers",
            "Content-Type, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods",
            "POST, GET, OPTIONS, PUT, DELETE")

        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
