package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Database      DatabaseConfig
	Storage       StorageConfig
	Log           LogConfig
	OpenAI_APIKey string `mapstructure:"openai_api_key"`
	TemplatesPath string `mapstructure:"TEMPLATES_PATH"`
	Image         ImageConfig
	Migration     MigrationConfig
	JWT           JWTConfig
}

type DatabaseConfig struct {
	Path string
}

type StorageConfig struct {
	UseS3    bool   `mapstructure:"useS3"`
	Region   string `mapstructure:"region"`
	Bucket   string `mapstructure:"bucket"`
	BasePath string `mapstructure:"basePath"`
}

type LogConfig struct {
	Dir  string
	File string
}

type ImageConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

func LoadConfig() (Configuration, error) {
	viper.SetConfigFile("./config/config.yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}

	var config Configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}

type MigrationConfig struct {
	Path string
}

type JWTConfig struct {
	Phrase string
}
