package download

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rjmateus/go-suma/config"
	"github.com/rjmateus/go-suma/repositories/download"
	"net/http"
	"os"
	"path"
	"strings"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		return true
	} else {
		return false
	}

}

const mountPoint = "/var/spacewalk/"

func HandleRepodata() gin.HandlerFunc {
	return func(c *gin.Context) {
		channel := c.Param("channel")
		fileName := c.Param("file")

		filePath := fmt.Sprintf("/var/cache/rhn/repodata/%s/%s", channel, fileName)
		if !PathExists(filePath) {
			if strings.HasSuffix(fileName, ".asc") || strings.HasSuffix(fileName, ".key") {
				c.String(http.StatusNotFound, fmt.Sprintf("Key or signature file not provided: %s", fileName))
			} else {
				c.String(http.StatusNotFound, fmt.Sprintf("File not found:%s", c.Request.URL.Path))
			}
		}
		downloadProcessor(c, filePath)
	}
}

func HandlePackage(app *config.Application) gin.HandlerFunc {

	return func(c *gin.Context) {
		channel := c.Param("channel")
		pkinfo := parsePackageFileName(c.Request.URL.Path)
		packageDb, error := download.GetDownloadPackage(app.DBGorm, channel, pkinfo.name, pkinfo.version, pkinfo.release, pkinfo.arch, pkinfo.checksum, pkinfo.epoch)
		if error != nil {
			c.String(http.StatusNotFound, fmt.Sprintf("%s not found", path.Base(c.Request.URL.Path)))
		}
		downloadProcessor(c, path.Join(mountPoint, packageDb.Path))
	}
}

func HandlerMediaFiles() gin.HandlerFunc {
	return func(c *gin.Context) {
		channelLabel := c.Param("channel")
		fileName := c.Param("file")

		if fileName == "products" {
			filePath := getMediaProductsFile(channelLabel)
			if len(filePath) == 0 {
				c.String(http.StatusNotFound, fmt.Sprintf("%s not found", fileName))
			} else {
				downloadProcessor(c, filePath)
			}
		} else {
			c.String(http.StatusNotFound, fmt.Sprintf("%s not found", fileName))
		}
	}
}

func getMediaProductsFile(channelLabel string) string {
	return ""
}

func downloadProcessor(c *gin.Context, filePath string) {
	if !PathExists(filePath) {
		c.String(http.StatusNotFound, fmt.Sprintf("File not found:%s", c.Request.URL.Path))
	} else {
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path.Base(filePath)))
		c.File(filePath)
	}
}

type pkgInfo struct {
	name     string
	version  string
	release  string
	epoch    string
	arch     string
	orgId    string
	checksum string
}

func parsePackageFileName(filepath string) pkgInfo {
	parts := strings.Split(filepath, "/")

	extension := path.Ext(filepath)
	basename := strings.TrimSuffix(path.Base(filepath), extension)
	arch := basename[strings.LastIndex(basename, ".")+1:]
	rest := basename[:strings.LastIndex(basename, ".")]

	//String arch = StringUtils.substringAfterLast(basename, ".");
	//String rest = StringUtils.substringBeforeLast(basename, ".");
	release := ""
	name := ""
	version := ""
	epoch := ""
	org := ""
	checksum := ""

	// Debian packages names need spacial handling
	if "deb" == extension || "udeb" == extension {
		name = rest[strings.LastIndex(rest, "_")+1:]
		rest = rest[:strings.LastIndex(rest, "-")]

		pkgEv := parseDebian(rest)
		epoch = pkgEv.epoch
		version = pkgEv.version
		release = pkgEv.release
	} else {
		release = rest[strings.LastIndex(rest, "-")+1:]
		rest = rest[:strings.LastIndex(rest, "-")]
		version = rest[strings.LastIndex(rest, "-")+1:]
		name = rest[:strings.LastIndex(rest, "-")]
	}
	// path is getPackage/<org>/<checksum>/filename
	if len(parts) == 9 && parts[5] == "getPackage" {
		org = parts[6]
		checksum = parts[7]
	}
	return pkgInfo{
		name:     name,
		version:  version,
		release:  release,
		epoch:    epoch,
		arch:     arch,
		orgId:    org,
		checksum: checksum,
	}
}

type debianPackage struct {
	epoch   string
	version string
	release string
}

func parseDebian(version string) debianPackage {

	// repo-sync replaces empty releases with 'X'. We copy the same behavior.
	release := "X"
	epoch := ""

	epochIndex := strings.Index(version, ":")
	if epochIndex > 0 {
		// Strip away optional 'epoch'
		epoch = version[:epochIndex]
		version = version[epochIndex+1:]
	}

	releaseIndex := strings.Index(version, "-")
	if releaseIndex > 0 {
		// Strip away optional 'release'
		release = version[releaseIndex+1:]
		version = version[:releaseIndex]
	}

	return debianPackage{epoch, version, release}
}
