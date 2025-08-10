package server

import (
	_ "github.com/example/donkey/docs" // swagger docs
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/example/donkey/internal/api"
	"github.com/example/donkey/internal/db"
	"github.com/example/donkey/internal/game"
	"github.com/example/donkey/internal/model"
)

// New creates a new HTTP server with routes configured.
func New() *gin.Engine {
	game.VerifyAssets()
	db.Init(&model.User{}, &model.Game{}, &model.GamePlayer{}, &model.GameCard{}, &model.GameState{})

	r := gin.Default()
	r.Static("/assets", "./web/assets")
	r.Static("/ui", "./web/ui")

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/hello", api.HelloHandler)
		apiGroup.POST("/register", api.RegisterHandler)
		apiGroup.POST("/game/start", api.StartGameHandler)
		apiGroup.POST("/game/join", api.JoinGameHandler)
		apiGroup.POST("/game/finalize", api.FinalizeHandler)
		apiGroup.GET("/game/state", api.GameStateHandler)
		apiGroup.GET("/version", api.VersionHandler)
		apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	return r
}
