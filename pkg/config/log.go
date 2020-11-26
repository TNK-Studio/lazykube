package config

import "github.com/sirupsen/logrus"

type LogConfig struct {
	Path  string       `yaml:"path"`
	Level logrus.Level `yaml:"level"`
}
