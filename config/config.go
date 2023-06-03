package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		AppName        string
		AppPort        int
		RedisConfig    CacheConfig
		MysqlConfig    DBConfig
		PostgresConfig DBConfig
	}

	CacheConfig struct {
		Host     string
		Port     int
		Password string
		DB       int
	}

	DBConfig struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
)

func New() *Config {
	return &Config{
		AppName: getEnv("APP_NAME", "go-boilerplate"),
		AppPort: getEnvAsInt("APP_PORT", 8080),
		RedisConfig: CacheConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		MysqlConfig: DBConfig{
			Host:     getEnv("DB_MYSQL_HOST", "localhost"),
			Port:     getEnvAsInt("DB_MYSQL_PORT", 3306),
			User:     getEnv("DB_MYSQL_USER", "mysql"),
			Password: getEnv("DB_MYSQL_PASSWORD", ""),
			DBName:   getEnv("DB_MYSQL_NAME", "db_solid"),
		},
		PostgresConfig: DBConfig{
			Host:     getEnv("DB_POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("DB_POSTGRES_PORT", 5432),
			User:     getEnv("DB_POSTGRES_USER", "postgres"),
			Password: getEnv("DB_POSTGRES_PASSWORD", ""),
			DBName:   getEnv("DB_POSTGRES_NAME", "db_solid"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}

	if nextValue := os.Getenv(key); nextValue != "" {
		return nextValue
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultVal
}
