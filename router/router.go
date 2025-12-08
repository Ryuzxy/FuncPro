package router

import (
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/ryuzxy/FuncPro/internal/middleware"
    "github.com/ryuzxy/FuncPro/pkg/komoditas"
    "github.com/ryuzxy/FuncPro/pkg/price"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // Middleware
    r.Use(middleware.Logger())
    r.Use(gin.Recovery()) // diganti karena middleware.Recovery() tidak ada
    r.Use(middleware.CORS())

    // Initialize repositories
    komoditasRepo := komoditas.NewRepository(db)
    priceRepo := price.NewPriceRepository(db)

    // Initialize services
    komoditasService := komoditas.NewService(komoditasRepo) // hanya 1 argumen
    priceService := price.NewService(priceRepo)

    // Initialize handlers
    komoditasHandler := komoditas.NewHandler(komoditasService)
    priceHandler := price.NewHandler(priceService)

    // API routes
    api := r.Group("/api/v1")
    {
        // Komoditas routes
        komoditasGroup := api.Group("/komoditas")
        {
            komoditasGroup.GET("", komoditasHandler.GetAllKomoditas)
            komoditasGroup.POST("", komoditasHandler.CreateKomoditas)
            komoditasGroup.GET("/:id", komoditasHandler.GetKomoditasByID)
            komoditasGroup.PUT("/:id", komoditasHandler.UpdateKomoditas)
            komoditasGroup.DELETE("/:id", komoditasHandler.DeleteKomoditas)
            komoditasGroup.GET("/:id/stats", komoditasHandler.GetKomoditasStats)
        }

        // Price routes
        priceGroup := api.Group("/prices")
        {
            priceGroup.POST("", priceHandler.CreatePrice)
            priceGroup.POST("/bulk", priceHandler.BulkCreatePrices)
            priceGroup.GET("/komoditas/:komoditas_id", priceHandler.GetPricesByKomoditas)
            priceGroup.GET("/komoditas/:komoditas_id/analysis", priceHandler.GetPriceAnalysis)
        }

        // Health check
        api.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "status":    "ok",
                "timestamp": time.Now().Unix(),
            })
        })
    }

    return r
}
