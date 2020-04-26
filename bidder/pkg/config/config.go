package config

import (
	"fmt"
	"time"

	"bidder/pkg/env"
)

const (
	appURLKey          = "APP_URL"
	appPortKey         = "APP_PORT"
	delayTimeKey       = "DELAY_TIME_MS"
	RegistrationURLKey = "REGISTRATION_URL"
)

type Config struct {
	AppURL          string
	AppPort         int
	Delay           time.Duration
	RegistrationURL string
}

func New() (*Config, error) {
	vars := &env.Vars{}

	appURL := vars.Mandatory(appURLKey)
	appPort := vars.MandatoryInt(appPortKey)
	delay := vars.MandatoryInt(delayTimeKey)
	registrationURL := vars.Mandatory(RegistrationURLKey)

	if err := vars.Error(); err != nil {
		return nil, fmt.Errorf("config: environment variables: %v", err)
	}

	config := &Config{
		AppURL:          appURL,
		AppPort:         appPort,
		Delay:           time.Duration(delay) * time.Millisecond,
		RegistrationURL: registrationURL,
	}
	return config, nil
}
