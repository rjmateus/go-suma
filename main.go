package main

import (
	"github.com/rjmateus/go-suma/config"
	"github.com/rjmateus/go-suma/web"
)

func main() {
	app := config.NewApplication()
	web.InitRoutes(app)
	app.Engine.Run(":8088")

	// https://hoohoo.top/blog/20210530112304-golang-tutorial-introduction-gin-html-template-and-how-integration-with-bootstrap/

}
