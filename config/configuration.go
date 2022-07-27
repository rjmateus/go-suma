package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func NewConfiguration() *SumaConfiguration {
	baseConfig := viper.New()
	baseConfig.SetConfigName("config") // name of config file (without extension)
	baseConfig.SetConfigFile("rhn.conf")
	baseConfig.SetConfigType("properties") // REQUIRED if the config file does not have the extension in the name
	baseConfig.AddConfigPath("/etc/rhn")   // path to look for the config file in
	baseConfig.AddConfigPath(".")          // optionally look for config in the working directory
	err := baseConfig.ReadInConfig()       // Find and read the config file
	if err != nil {                        // Handle errors reading the config file
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
