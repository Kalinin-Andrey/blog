package integration

import (
	"errors"
)

type AppConfig struct {
	NameSpace   string
	Subsystem   string
	Service     string
	Environment string
}

type Integration struct {
}

func New(appConfig *AppConfig, cfg *Config) (*Integration, error) {
	integration := &Integration{}

	return integration, nil
}

func (intgr *Integration) Close() error {
	return errors.Join()
}
