package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var defaults = map[string]interface{}{
	"mountPoint":                      "/var/spacewalk/",
	"java.salt_check_download_tokens": "true",
}

func NewConfiguration() *SumaConfiguration {
	baseConfig := viper.New()
	baseConfig.SetConfigName("config") // name of config file (without extension)
	baseConfig.SetConfigFile("/etc/rhn/rhn.conf")
	baseConfig.SetConfigType("properties") // REQUIRED if the config file does not have the extension in the name
	baseConfig.AddConfigPath("/etc/rhn/")  // path to look for the config file in
	//baseConfig.AddConfigPath(".")          // optionally look for config in the working directory

	for key, value := range defaults {
		baseConfig.SetDefault(key, value)
	}

	baseConfig.SetEnvPrefix("suma")
	baseConfig.AutomaticEnv()

	err := baseConfig.ReadInConfig() // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Println(baseConfig.AllKeys())

	return &SumaConfiguration{baseConfig: baseConfig}
}

type SumaConfiguration struct {
	baseConfig *viper.Viper
}

func (c *SumaConfiguration) GetMountPoint() string {
	return c.baseConfig.GetString("mountPoint")
}

func (c *SumaConfiguration) GetString(key string) string {
	return c.baseConfig.GetString(key)
}

func (c *SumaConfiguration) CheckDownloadToken() bool {
	return c.baseConfig.GetBool("java.salt_check_download_tokens")
}

func (c *SumaConfiguration) GetServerSecretKey() string {
	return c.baseConfig.GetString("server.secret_key")
}
