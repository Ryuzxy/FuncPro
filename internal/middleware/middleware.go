package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/ryuzxy/fucpro/pkg/fx"
)

// Logger middleware
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        fx.Try(func() (interface{}, error) {
            gin.Default().Errorf("[%s] %s %s %d %v",
                c.Request.Method,
                c.Request.URL.Path,
                c.ClientIP(),
                c.Writer.Status(),
                duration,
            )
            return nil, nil
        })
    }
}

// Recovery middleware
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                c.AbortWithStatusJSON(500, gin.H{
                    "success": false,
                    "error":   "Internal server error",
                })
            }
        }()
        c.Next()
    }
}

// CORS middleware
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}