package download

import (
	"github.com/rjmateus/go-suma/config"
	"github.com/rjmateus/go-suma/handlers/download"
)

func InitDownloadRoutes(app *config.Application) {

	router := app.Engine.Group("/rhn/manager/download")
	router.GET("/:channel/repodata/:file", download.HandleRepodata())
	router.HEAD("/:channel/repodata/:file", download.HandleRepodata())

	router.GET("/:channel/media.1/:file", download.HandlerMediaFiles())
	router.HEAD("/:channel/media.1/:file", download.HandlerMediaFiles())

	router.GET("/:channel/getPackage/:org/:checksum/:file", download.HandlePackage(app))
	router.HEAD("/:channel/getPackage/:org/:checksum/:file", download.HandlePackage(app))
	// :org doesn't represent org, in this constext it represents the file.
	// this is a known limitation of gin. See:
	// - https://github.com/gin-gonic/gin/issues/1301#issuecomment-392346179
	// - https://github.com/gin-gonic/gin/issues/1681
	router.GET("/:channel/getPackage/:org", download.HandlePackage(app))
	router.HEAD("/:channel/getPackage/:org", download.HandlePackage(app))
}
