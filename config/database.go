package config

import (
	"fmt"
	_ "github.com/lib/pq"
)

// getConnectionString return the connection string for the database after reading config file for
func getConnectionString(configFilePath *SumaConfiguration) string {

	host := configFilePath.GetString("db_host")
	port := configFilePath.GetString("db_port")
	dbname := configFilePath.GetString("db_name")
	user := configFilePath.GetString("db_user")
	password := configFilePath.GetString("db_password")

	return fmt.Sprintf("user='%s' password='%s' dbname='%s' host='%s' port='%s' sslmode=disable",
		user, password, dbname, host, port)
}
