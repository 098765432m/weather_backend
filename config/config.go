package config

import (
	"github.com/098765432m/logger"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	ParseTime bool  `mapstructure:"parse_time"`
}

type WeatherApiConfig struct {
	ApiKey string `mapstructure:"api_key"`
}

type JWTConfig struct {
	SecretKey string `mapstructure:"secret_key"`
}

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Database   DatabaseConfig   `mapstructure:"database"`
	WeatherApp WeatherApiConfig `mapstructure:"weather_api"`
	JWT		   JWTConfig     	`mapstructure:"jwt"`
}

var AppData Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config") // Fallback to this path if path above fails

	if err := viper.ReadInConfig(); err != nil {
		logger.NewLogger().Error.Fatal("Cannot read config file: ", err)
	}

	if err := viper.Unmarshal(&AppData); err != nil {
		logger.NewLogger().Error.Fatal("Cannot unmarshal config file: ", err)

	}
}
