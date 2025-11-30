package main

import (
    "log"
    "os"
    
    "github.com/joho/godotenv"
    "github.com/ryuzxy/fucpro/db"
    "github.com/ryuzxy/fucpro/internal/config"
    "github.com/ryuzxy/fucpro/router"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    // Load configuration
    cfg := config.Load()
    
    // Initialize database
    database, err := db.InitDB(cfg)
    if err != nil {
        log.Printf("Failed to initialize database: %v", err)
        os.Exit(1)
    }
    
    // Setup router
    r := router.SetupRouter(database)
    
    // Start server
    log.Printf("Server starting on port %s", cfg.ServerPort)
    if err := r.Run(":" + cfg.ServerPort); err != nil {
        log.Printf("Server failed: %v", err)
        os.Exit(1)
    }
}