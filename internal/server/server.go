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
	db.Init(&model.User{}, &model.Game{}, &model.GamePlayer{}, &model.GameCard{}, &model.GameState{}, &model.GameSessionLog{})

	r := gin.Default()
	r.Use(logRequests())
	r.Static("/assets", "./web/assets")
	r.Static("/ui", "./web/ui")
	r.StaticFile("/favicon.ico", "./web/assets/favicon.ico")
	r.StaticFile("/apple-touch-icon.png", "./web/assets/apple-touch-icon.png")
	r.StaticFile("/apple-touch-icon-precomposed.png", "./web/assets/apple-touch-icon.png")

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/register", api.RegisterHandler)
		apiGroup.POST("/game/start", api.StartGameHandler)
		apiGroup.POST("/game/join", api.JoinGameHandler)
		apiGroup.POST("/game/finalize", api.FinalizeHandler)
		apiGroup.POST("/game/abandon", api.AbandonHandler)
		apiGroup.POST("/game/chat", api.ChatHandler)
		apiGroup.GET("/game/:gameId/logs", api.LogsHandler)
		apiGroup.GET("/game/:gameId/stream/:userId", api.StreamHandler)
		apiGroup.GET("/game/:gameId/state/:userId", api.GameStateHandler)
		apiGroup.GET("/admin/game/:gameId/state", api.AdminStateHandler)
		apiGroup.POST("/user/:id/rename", api.RenameHandler)
		apiGroup.GET("/user/:id", api.GetUserHandler)
		apiGroup.GET("/users", api.ListUsersHandler)
		apiGroup.GET("/version", api.VersionHandler)
		apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
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
