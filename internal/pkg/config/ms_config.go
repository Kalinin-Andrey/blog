package config

import (
	"fmt"

	"github.com/spf13/viper"

	"blog/internal/integration"
)

type MSConfiguration struct {
	App         *AppConfig
	Integration *integration.UsageConfig
}

func (c *Configuration) readMSConfig(pathToConfig string) error {
	msConfig := &MSConfiguration{}

	viper.SetConfigFile(pathToConfig)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("MSConfig file not found in %q", pathToConfig)
		} else {
			return fmt.Errorf("MSConfig file was found in %q, but was produced error: %w", pathToConfig, err)
		}
	}

	err := viper.Unmarshal(msConfig)
	if err != nil {
		return fmt.Errorf("MSConfig unmarshal error: %w", err)
	}

	if c.applyMSConfig(msConfig); err != nil {
		return fmt.Errorf("applyMSConfig error: %w", err)
	}
	return nil
}

func (c *Configuration) applyMSConfig(config *MSConfiguration) error {
	if config.App.Environment != "" {
		c.App.Environment = config.App.Environment
	}
	if config.App.NameSpace != "" {
		c.App.NameSpace = config.App.NameSpace
	}
	if config.App.Name != "" {
		c.App.Name = config.App.Name
	}
	if config.App.Service != "" {
		c.App.Service = config.App.Service
	}
	return nil
}
