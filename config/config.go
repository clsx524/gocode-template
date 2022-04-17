package config

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/config"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// EnvironmentKey environment variables storing the current deployment environment
	EnvironmentKey = "RENV"
	// EnvTest test env
	EnvTest = "test"
	// EnvProduction production env
	EnvProduction = "production"
	// EnvStaging staging env
	EnvStaging = "staging"

	baseFile         = "base"
	defaultConfigDir = "config"
)

type Provider interface {
	LoggerConfig() zap.Config
	ServiceConfig() ServiceConfig
	MongoConfig() MongoConfig
}

type AppConfig struct {
	Logger  zap.Config
	Service ServiceConfig `yaml:"service"`
	MongoDB MongoConfig   `yaml:"mongodb"`
}

type ServiceConfig struct {
	Name        string `yaml:"name"`
	Environment string
	Port        string `yaml:"port"`
}

type MongoConfig struct {
	URI string `yaml:"uri"`
}

func (c *AppConfig) LoggerConfig() zap.Config {
	return c.Logger
}

func (c *AppConfig) ServiceConfig() ServiceConfig {
	return c.Service
}

func (c *AppConfig) MongoConfig() MongoConfig {
	return c.MongoDB
}

func New() (Provider, error) {
	var c AppConfig
	provider, err := getConfigProvider()
	if err != nil {
		return nil, err
	}

	if err := provider.Get("").Populate(&c); err != nil {
		return nil, err
	}

	c.Service.Environment = getEnvironment()

	if c.Service.Environment == EnvTest {
		c.Logger = zap.NewDevelopmentConfig()
	} else {
		c.Logger = zap.NewProductionConfig()
	}

	return &c, nil
}

func getEnvironment() string {
	env := os.Getenv(EnvironmentKey)
	if env == "" {
		return EnvTest
	}

	return env
}

func getConfigProvider() (config.Provider, error) {
	names := []string{baseFile, getEnvironment()}

	fileNames := make([]string, len(names))
	for idx, name := range names {
		fileNames[idx] = filepath.Join(defaultConfigDir, fmt.Sprintf("%s.yaml", name))
	}

	found := false
	opts := make([]config.YAMLOption, 0)
	for _, info := range fileNames {
		contents, err := ioutil.ReadFile(info)
		if err != nil && os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, err
		}

		found = true
		r := bytes.NewReader(contents)
		opts = append(opts, config.RawSource(r))
	}

	if !found {
		return nil, errors.New("failed to find configuration files")
	}

	return config.NewYAML(opts...)
}
