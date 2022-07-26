package web

import (
	"github.com/gin-gonic/gin"
	"github.com/rjmateus/go-suma/config"
	"github.com/rjmateus/go-suma/web/download"
	"net/http"
)

func initLocal(engine *gin.Engine) {
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong\n")
	})
}

func InitRoutes(app *config.Application) {
	initLocal(app.Engine)
	download.InitDownloadRoutes(app)
}
