package config

import (
	"fmt"

	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/source/cli"
)

type Settings struct {
	Tube     string
	Priority uint32
	Delay    int
	TTR      int
}

type Config struct {
	Host     string
	Port     int
	settings *Settings
}

func (c Config) Current() string {

	return fmt.Sprintf(("------------------------\n" +
		"Address: %s:%d\n" +
		"Tube: %s\n" +
		"Priority: %d\n" +
		"Delay: %d\n" +
		"Time to run: %d\n" +
		"------------------------\n"), c.Host, c.Port, c.settings.Tube, c.settings.Priority, c.settings.Delay, c.settings.TTR)
}

func (c *Config) GetSettings() *Settings {
	return c.settings
}

func (c *Config) ChangeTube(tube string) {
	c.settings.Tube = tube
}

func (c *Config) ChangeDelay(delay int) {
	c.settings.Delay = delay
}

func (c *Config) ChangeTTR(ttr int) {
	c.settings.TTR = ttr
}

func (c *Config) ChangePriority(priority uint32) {
	c.settings.Priority = priority
}

var defaultConfig = Config{
	settings: &Settings{},
}

func init() {
	defaultConfig.settings.Tube = "default"
	defaultConfig.settings.Priority = 1
	defaultConfig.settings.TTR = 60
}

func Load() (*Config, error) {
	var cfg = defaultConfig
	configor, err := config.NewConfig(
		config.WithSource(cli.NewSource()),
	)
	if err != nil {
		return nil, fmt.Errorf("error configure config source. %w", err)
	}

	if err := configor.Load(); err != nil {
		return nil, fmt.Errorf("error load config. %w", err)
	}

	if err := configor.Scan(&cfg); err != nil {
		return nil, fmt.Errorf("error scan config. %w", err)
	}
	return &cfg, nil
}
