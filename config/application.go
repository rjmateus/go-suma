package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	Engine *gin.Engine
	DBGorm *gorm.DB
}

var conf = "/etc/rhn/rhn.conf"

//var conf = "rhn.conf"

func NewApplication() *Application {
	r := gin.Default()
	dsn := getConnectionString(conf)
	dbGorm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Application{Engine: r, DBGorm: dbGorm}
}
