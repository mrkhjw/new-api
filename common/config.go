package common

import (
	"os"
	"strconv"
	"sync"
)

// Version information
var (
	Version   = "v0.0.1"
	BuildTime = "unknown"
)

// Server configuration
var (
	ServerAddress = "0.0.0.0"
	ServerPort    = 3000
	GinMode       = "release" // changed from "debug" - safer default for personal deployment
)

// Database configuration
var (
	SQLitePath = "new-api.db"
	DBType     = "sqlite" // sqlite, mysql, postgres
	DBHost     = "localhost"
	DBPort     = 3306
	DBName     = "new-api"
	DBUser     = "root"
	DBPassword = ""
)

// Redis configuration
var (
	RedisConnString = ""
	RedisPassword   = ""
	RedisDB         = 0
	UsingRedis      = false
)

// Security configuration
var (
	SessionSecret = "new-api-secret"
	CryptoSecret  = ""
	RootUserEmail = ""
	RootUserName  = "root"
	RootUserPwd   = "123456"
)

// System configuration
var (
	SystemName      = "New API"
	SystemLogo      = ""
	FooterHTML      = ""
	HomePageContent = ""
	Theme           = "default"
	EnableSwagger   = false
	DebugEnabled    = false
)

// Rate limiting
var (
	GlobalApiRateLimitNum      = 60
	GlobalApiRateLimitDuration = int64(3 * 60)
)

// Billing configuration
var (
	QuotaPerUnit     = 500 * 1000.0 // $1 = 500k tokens by default
	InitialRootToken = ""
)

var mu sync.RWMutex

// LoadConfigFromEnv loads configuration values from environment variables,
// falling back to defaults if not set.
func LoadConfigFromEnv() {
	mu.Lock()
	defer mu.Unlock()

	if v := os.Getenv("SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			ServerPort = port
		}
	}
	if v := os.Getenv("GIN_MODE"); v != "" {
		GinMode = v
	}
	if v := os.Getenv("SQL_DSN"); v != "" {
		// If SQL_DSN is provided, use it directly (handled by database package)
	}
	if v := os.Getenv("REDIS_CONN_STRING"); v != "" {
		RedisConnString = v
		UsingRedis = true
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		RedisPassword = v
	}
	if v := os.Getenv("SESSION_SECRET"); v != "" {
		SessionSecret = v
	}
	if v := os.Getenv("CRYPTO_SECRET"); v != "" {
		CryptoSecret = v
	}
	if v := os.Getenv("SYSTEM_NAME"); v != "" {
		SystemName = v
	}
	if v := os.Getenv("ROOT_USER_EMAIL"); v != "" {
		RootUserEmail = v
	}
	if v := os.Getenv("INITIAL_ROOT_TOKEN"); v != "" {
		InitialRootToken = v
	}
	if v := os.Getenv("DEBUG"); v == "true" || v == "1" {
		DebugEnabled = true
	}
}
