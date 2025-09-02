package server

import (
	"log"
	"net/http"

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
	// Serve static assets for both development and production
	r.Static("/assets", "./dist/assets")           // Vue.js build assets (production)
	r.Static("/web/assets", "./web/assets")        // Game assets (card images, icons)
	r.StaticFile("/favicon.ico", "./web/assets/favicon.ico")
	r.StaticFile("/apple-touch-icon.png", "./web/assets/apple-touch-icon.png")
	r.StaticFile("/apple-touch-icon-precomposed.png", "./web/assets/apple-touch-icon.png")

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

	// Serve Vue.js app
	r.GET("/", func(c *gin.Context) {
		c.File("./dist/index.html")
	})
	
	// Fallback for SPA routing
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

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
