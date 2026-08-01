package configs

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

func GetEnv(key string, value ...string) string {
	if err := godotenv.Load(".env"); err != nil {
		// Try loading from workspace root or system environment if .env load fails
		_ = godotenv.Load()
	}

	if os.Getenv(key) != "" {
		return os.Getenv(key)
	} else {
		if len(value) > 0 {
			return value[0]
		}
		return ""
	}
}

var (
	APP_PORT             = GetEnv("APP_PORT", "10002")
	APP_KEY              = GetEnv("APP_KEY", "SomeSecretRandomKeyHereForJWTSigning")
	APP_ENV              = GetEnv("APP_ENV", "DEV")
	DB_DRIVER            = GetEnv("DB_DRIVER", "postgres")
	DB_HOST              = GetEnv("DB_HOST", "localhost")
	DB_PORT              = GetEnv("DB_PORT", "5432")
	DB_NAME              = GetEnv("DB_NAME", "game_biller")
	DB_USER              = GetEnv("DB_USER", "postgres")
	DB_PASS              = GetEnv("DB_PASS", "1234")
	SSL_MODE             = GetEnv("SSL_MODE", "disable")
	HOST_WORKER          = GetEnv("HOST_WORKER")
	ENDPOINT_IAK_PAYMENT = GetEnv("ENDPOINT_IAK_PAYMENT")
	ENDPOINT_IAK_INQUIRY = GetEnv("ENDPOINT_IAK_INQUIRY")
	DEFAULT_MID          = GetEnv("DEFAULT_MID")
)

func GetLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.UTC
	}
	return loc
}
