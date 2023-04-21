package config

import (
	"log"

	"github.com/spf13/viper"
)

// ServerConfig struct to handle server configuration
type ServerConfig struct {
	Addr            string
	WriteTimeout    int
	ReadTimeout     int
	GraceFulTimeout int
}

// DBConfig struct to handle database configuration
type DBConfig struct {
	Name            string
	Host            string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime int
}

// Config struct for .env.yml
type Config struct {
	Server ServerConfig
	DB     DBConfig
}

// InitConfig function to init configuration, returns Config struct
func InitConfig() Config {
	viper.SetConfigName(".env")

	viper.AddConfigPath(".")

	var configuration Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return configuration
}
