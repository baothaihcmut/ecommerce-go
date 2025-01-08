package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Consol   ConsulConfig
	Jwt      JwtConfig
}

type JwtConfig struct {
	AccessTokenSecret  string
	AccessTokenAge     int
	RefreshTokenSecret string
	RefreshTokenAge    int
}

type ServerConfig struct {
	Host              string
	Port              int
	MaxConnectionIdle int
	MaxConnectionAge  int
	Time              int
	Timeout           int
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type ConsulConfig struct {
	Host string
	Port int
}

func getEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Could not parse environment variable %s as int, using fallback %d", key, fallback)
	}
	return fallback
}

//	func getEnvBool(key string, fallback bool) bool {
//		if value, exists := os.LookupEnv(key); exists {
//			if boolValue, err := strconv.ParseBool(value); err == nil {
//				return boolValue
//			}
//			log.Printf("Warning: Could not parse environment variable %s as bool, using fallback %t", key, fallback)
//		}
//		return fallback
//	}
func LoadConfig(env string) (*Config, error) {

	if env == "dev" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}
	return &Config{
		Server: ServerConfig{
			Host:              getEnvString("SERVER_HOST", "localhost"),
			Port:              getEnvInt("SERVER_PORT", 8080),
			MaxConnectionIdle: getEnvInt("MAX_CONNECTION_IDLE", 5),
			MaxConnectionAge:  getEnvInt("MAX_CONNECTION_AGE", 10),
			Time:              getEnvInt("TIME", 2),
			Timeout:           getEnvInt("TIMEOUT", 20),
		},
		Database: DatabaseConfig{
			Driver:   getEnvString("DB_DRIVER", "postgres"),
			Host:     getEnvString("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnvString("DB_USER", "postgres"),
			Password: getEnvString("DB_PASSWORD", "postgres"),
			Name:     getEnvString("DB_NAME", "user"),
		},
		Consol: ConsulConfig{
			Host: getEnvString("CONSUL_HOST", "localhost"),
			Port: getEnvInt("CONSUL_PORT", 8500),
		},
		Jwt: JwtConfig{
			AccessTokenSecret:  getEnvString("JWT_ACCESS_TOKEN_SECRET", "access_token_access_secret"),
			AccessTokenAge:     getEnvInt("JWT_ACCESS_TOKEN_AGE", 3),
			RefreshTokenSecret: getEnvString("JWT_REFRESH_TOKEN_SECRET", "refresh_token_access_secret"),
			RefreshTokenAge:    getEnvInt("JWT_REFRESH_TOKEN_AGE", 7),
		},
	}, nil

}
