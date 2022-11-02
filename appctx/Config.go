package appctx

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Config struct {
	DbDns             string
	DbMaxIdleConns    int
	DbMaxOpenConns    int
	DbConnMaxLifetime int // seconds
	ServerAddr        string
}

// CreateConfigFromEnv 환경변수로부터 Config 생성
func CreateConfigFromEnv() Config {
	var config = Config{}
	config.DbDns = os.Getenv("DB_DNS")
	config.DbMaxIdleConns = getOrDefaultInt("DB_MAX_IDLE_CONNS", 10)
	config.DbMaxOpenConns = getOrDefaultInt("DB_MAX_OPEN_CONNS", 40)
	config.DbConnMaxLifetime = getOrDefaultInt("DB_CONN_MAX_LIFETIME", 60)
	config.ServerAddr = getOrDefaultString("SERVER_ADDR", ":8080")
	return config
}

func getOrDefaultString(env string, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		log.Warnf("env[%s]'s value is empty, so use default %s", env, defaultValue)
		return defaultValue
	}
	return value
}

func getOrDefaultInt(env string, defaultValue int) int {
	value := os.Getenv(env)
	if value == "" {
		log.Warnf("env[%s]'s value is empty, so use default %d", env, defaultValue)
		return defaultValue
	}
	intValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		log.Warnf("env[%s]'s value %s is not number, so use default %d", env, value, defaultValue)
		return defaultValue
	}
	return int(intValue)
}
