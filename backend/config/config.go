package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Google   *GoogleConfig
	Server   *ServerConfig
	Secrets  *SecretConfig
	Frontend *FrontendConfig
}

type ServerConfig struct {
	Secure bool
	Domain string
	Base   string
	Host   string
	Port   int64
}

type GoogleConfig struct {
	ClientId     string
	ClientSecret string
}

type SecretConfig struct {
	Jwt string
}

type FrontendConfig struct {
	Base string
}

var config *Config

func InitConfig(name string) *Config {
	viper.AddConfigPath("config")
	viper.SetConfigName(name)
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error while reading config: %s", err.Error()))
	}

	serverConfig := ServerConfig{
		Secure: viper.Get("server.secure").(bool),
		Domain: viper.Get("server.domain").(string),
		Base:   viper.Get("server.base").(string),
		Host:   viper.Get("server.host").(string),
		Port:   viper.Get("server.port").(int64),
	}

	googleConfig := GoogleConfig{
		ClientId:     viper.Get("google.client_id").(string),
		ClientSecret: viper.Get("google.client_secret").(string),
	}

	secretConfig := SecretConfig{
		Jwt: viper.Get("secrets.jwt").(string),
	}

	frontendConfig := FrontendConfig{
		Base: viper.Get("frontend.base").(string),
	}

	config = &Config{
		Server:   &serverConfig,
		Google:   &googleConfig,
		Secrets:  &secretConfig,
		Frontend: &frontendConfig,
	}

	return config
}

func GetConfig() *Config {
	if config != nil {
		return config
	}
	return InitConfig("config")
}
