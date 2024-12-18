package config

import (
	"encoding/json"
	"sync/atomic"

	"github.com/spf13/viper"
)

var confValue atomic.Value

// Config stores all configuration of the application.
type info struct {
	Base struct {
		SourceID     int    `mapstructure:"source-id"`
		DownloadPath string `mapstructure:"download-path"`
		Extname      string `mapstructure:"extname"`
		AutoUpdate   int    `mapstructure:"auto-update"`
		LogLevel     string `mapstructure:"log-level"`
	} `mapstructure:"base"`
	Crawl struct {
		Threads int `mapstructure:"threads"`
	} `mapstructure:"crawl"`
	Retry struct {
		MaxAttempts int `mapstructure:"max-attempts"`
	} `mapstructure:"retry"`
}

func init() {
	confValue.Store(info{})
}

// ToJSON returns the JSON string representation of the Config
func (i info) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) error {
	viper.Reset()

	if path != "" {
		// Use the provided config file
		viper.SetConfigFile(path)
	} else {
		// If no config file is provided, use default locations
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var newConf info
	err = viper.Unmarshal(&newConf)
	if err != nil {
		return err
	}

	confValue.Store(newConf)
	return nil
}

func GetConf() info {
	return confValue.Load().(info)
}
