package server

import (
	_ "github.com/example/donkey/docs" // swagger docs
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/example/donkey/internal/api"
)

// New creates a new HTTP server with routes configured.
func New() *gin.Engine {
	r := gin.Default()

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/hello", api.HelloHandler)
		apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	return r
}
