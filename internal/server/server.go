package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/kairodrad/donkey/docs" // swagger docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/kairodrad/donkey/internal/api"
	"github.com/kairodrad/donkey/internal/db"
	"github.com/kairodrad/donkey/internal/game"
	"github.com/kairodrad/donkey/internal/model"
)

// New creates a new HTTP server with routes configured.
func New() *gin.Engine {
	game.VerifyAssets()
	db.Init(&model.User{}, &model.Game{}, &model.GamePlayer{}, &model.Round{}, &model.RoundPlayer{}, &model.Turn{}, &model.Card{}, &model.PlayedCard{}, &model.BotMemory{}, &model.GameSessionLog{}, &model.GameSettings{})

	// Set up publishers for game events
	game.SetStatePublisher(api.PublishState)
	game.SetLogPublisher(api.PublishLog)

	r := gin.Default()
	r.Use(logRequests())
	
	// Environment detection - default to production if dist/index.html exists
	isDev := os.Getenv("NODE_ENV") != "production"
	
	// Override: if dist/index.html exists and NODE_ENV is not explicitly set to dev, use production mode
	if _, err := os.Stat("./dist/index.html"); err == nil {
		if os.Getenv("NODE_ENV") == "" || os.Getenv("NODE_ENV") == "production" {
			isDev = false
		}
	}
	
	// Static file serving - works for both dev and production
	// Always serve game assets from /web/assets (consistent path)
	r.Static("/web/assets", "./web/assets")
	
	// In production, also serve built Vue.js assets
	if !isDev {
		r.Static("/assets", "./dist/assets")
		r.StaticFile("/favicon.ico", "./dist/assets/favicon-92640d13.ico")
		r.StaticFile("/apple-touch-icon.png", "./dist/assets/apple-touch-icon-71cd5b1a.png")
		r.StaticFile("/apple-touch-icon-precomposed.png", "./dist/assets/apple-touch-icon-71cd5b1a.png")
	}

	apiGroup := r.Group("/api")
	{
		// User management
		apiGroup.POST("/register", api.RegisterHandler)
		apiGroup.GET("/user/:id", api.GetUserHandler)
		apiGroup.GET("/users", api.ListUsersHandler)
		
		// Game management
		apiGroup.POST("/game/create", api.CreateGameHandler)
		apiGroup.POST("/game/join", api.JoinGameHandler)
		apiGroup.POST("/game/start", api.StartGameHandler)
		apiGroup.POST("/game/abandon", api.AbandonGameHandler)
		apiGroup.POST("/game/add-bot", api.AddBotHandler)
		apiGroup.POST("/game/play-card", api.PlayCardHandler)
		apiGroup.GET("/game/:gameId/state/:userId", api.GameStateHandler)
		apiGroup.GET("/games", api.GetGameListHandler)
		
		// Admin endpoints
		apiGroup.GET("/admin/game/:gameId/state", api.AdminStateHandler)
		
		// Chat and streaming
		apiGroup.POST("/game/chat", api.ChatHandler)
		apiGroup.GET("/game/:gameId/logs", api.LogsHandler)
		apiGroup.GET("/game/:gameId/stream/:userId", api.StreamHandler)
		
		// Utilities
		apiGroup.GET("/version", api.VersionHandler)
		apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Serve Vue.js app (only in production - dev uses Vite server)
	if !isDev {
		r.GET("/", func(c *gin.Context) {
			c.File("./dist/index.html")
		})
		
		// SPA fallback for production
		r.NoRoute(func(c *gin.Context) {
			// Don't serve SPA for API routes
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
				return
			}
			
			// Serve Vue.js app for client-side routing
			c.File("./dist/index.html")
		})
	} else {
		// In development, return a simple message for root
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Donkey API Server (Development Mode)",
				"frontend": "Run 'npm run dev' for Vue.js development server",
				"api": "/api/*",
				"assets": "/web/assets/*",
			})
		})
	}

	return r
}

// logRequests logs basic request info for debugging.
func logRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		if m := c.Request.Method; m == http.MethodGet || m == http.MethodPost || m == http.MethodPut {
			log.Printf("%s %s", m, c.Request.URL.Path)
		}
		c.Next()
	}
}
