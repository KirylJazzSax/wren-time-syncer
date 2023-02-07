package utils

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Username       string `mapstructure:"USERNAME"`
	JiraToken      string `mapstructure:"JIRA_TOKEN"`
	LinkAuthHeader string `mapstructure:"LINK_AUTH_HEADER"`
	JiraHost       string `mapstructure:"JIRA_HOST"`
	LinkHost       string `mapstructure:"LINK_HOST"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
