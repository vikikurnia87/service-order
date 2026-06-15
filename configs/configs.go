package configs

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	AppName            string
	ServerEnv          string
	ServerName         string
	ServerURL          string
	ServerPort         string
	GrpcPort           string
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ServerIdleTimeout  time.Duration
	ServerLang         string

	// UserGrpcAddr = alamat gRPC service-user (untuk ValidateToken/GetUser).
	UserGrpcAddr string

	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseNameTest string

	DatabaseTimeout      time.Duration
	DatabaseDialTimeout  time.Duration
	DatabaseReadTimeout  time.Duration
	DatabaseWriteTimeout time.Duration

	DatabaseSSLMode     string
	DatabaseMaxConn     int
	DatabaseMaxIdleConn int
	DatabaseMaxLifeConn int

	RedisNamespace        string
	RedisHost             string
	RedisPort             string
	RedisPassword         string
	RedisExpirationS      time.Duration
	RedisPaginationTTLS   time.Duration
	RedisPoolSize         int
	RedisMinIdleConns     int
	RedisDialTimeoutS     time.Duration
	RedisReadTimeoutS     time.Duration
	RedisWriteTimeoutS    time.Duration
	RedisPoolTimeoutS     time.Duration
	RedisMaxRetries       int
	RedisMinRetryBackoffS time.Duration
	RedisMaxRetryBackoffS time.Duration
	RedisIdleTimeoutS     time.Duration
	RedisMaxConnAgeS      time.Duration

	ElasticAPMServerURL      string
	ElasticAPMServiceName    string
	ElasticAPMServiceVersion string
	ElasticAPMEnvironment    string
)

func LoadEnv() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		slog.Error("cannot resolve caller path")
		return
	}
	root := filepath.Dir(filepath.Dir(filename))
	envPath := filepath.Join(root, ".env")
	if err := godotenv.Load(envPath); err != nil {
		slog.Error("Error loading .env file", "path", envPath, "error", err)
	}

	AppName = GetEnv("APP_NAME", "AMANOS-Vendor")
	ServerEnv = GetEnv("SERVER_ENV", "development")
	ServerName = GetEnv("SERVER_NAME", "service-order")
	ServerURL = GetEnv("SERVER_URL", "127.0.0.1:6004")
	ServerPort = GetEnv("SERVER_PORT", "6006")
	GrpcPort = GetEnv("GRPC_PORT", "61006")
	ServerReadTimeout = time.Duration(GetEnvAsInt("SERVER_READ_TIMEOUT", 10)) * time.Second
	ServerWriteTimeout = time.Duration(GetEnvAsInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second
	ServerIdleTimeout = time.Duration(GetEnvAsInt("SERVER_IDLE_TIMEOUT", 15)) * time.Second
	ServerLang = GetEnv("SERVICE_LANG", "en")

	UserGrpcAddr = GetEnv("USER_GRPC_ADDR", "127.0.0.1:6101")

	DatabaseHost = GetEnv("DB_HOST", "localhost")
	DatabasePort = GetEnv("DB_PORT", "5432")
	DatabaseUser = GetEnv("DB_USER", "postgres")
	DatabasePassword = GetEnv("DB_PASSWORD", "password")
	DatabaseName = GetEnv("DB_NAME", "service_order")
	DatabaseNameTest = GetEnv("DB_NAME_TEST", "service_order_test")

	DatabaseTimeout = time.Duration(GetEnvAsInt("DB_TIMEOUT", 30)) * time.Second
	DatabaseDialTimeout = time.Duration(GetEnvAsInt("DB_DIAL_TIMEOUT", 10)) * time.Second
	DatabaseReadTimeout = time.Duration(GetEnvAsInt("DB_READ_TIMEOUT", 30)) * time.Second
	DatabaseWriteTimeout = time.Duration(GetEnvAsInt("DB_WRITE_TIMEOUT", 15)) * time.Second

	DatabaseSSLMode = GetEnv("DB_SSLMODE", "disable")
	DatabaseMaxConn = GetEnvAsInt("DB_MAX_CONNECTIONS", 100)
	DatabaseMaxIdleConn = GetEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 10)
	DatabaseMaxLifeConn = GetEnvAsInt("DB_MAX_LIFETIME_CONNECTIONS", 2)

	RedisNamespace = GetEnv("REDIS_NAMESPACE", "service-order")
	RedisHost = GetEnv("REDIS_HOST", "")
	RedisPort = GetEnv("REDIS_PORT", "6379")
	RedisPassword = GetEnv("REDIS_PASSWORD", "")
	RedisExpirationS = time.Duration(GetEnvAsInt("REDIS_EXPIRATION_S", 900)) * time.Second
	RedisPaginationTTLS = time.Duration(GetEnvAsInt("REDIS_PAGINATION_TTL_S", 2700)) * time.Second
	RedisPoolSize = GetEnvAsInt("REDIS_POOLSIZE", 20)
	RedisMinIdleConns = GetEnvAsInt("REDIS_MIN_IDLE_CONNS", 4)
	RedisDialTimeoutS = time.Duration(GetEnvAsInt("REDIS_DIAL_TIMEOUT_S", 2)) * time.Second
	RedisReadTimeoutS = time.Duration(GetEnvAsInt("REDIS_READ_TIMEOUT_S", 2)) * time.Second
	RedisWriteTimeoutS = time.Duration(GetEnvAsInt("REDIS_WRITE_TIMEOUT_S", 2)) * time.Second
	RedisPoolTimeoutS = time.Duration(GetEnvAsInt("REDIS_POOL_TIMEOUT_S", 1)) * time.Second
	RedisMaxRetries = GetEnvAsInt("REDIS_MAX_RETRIES", 2)
	RedisMinRetryBackoffS = time.Duration(GetEnvAsInt("REDIS_MIN_RETRY_BACKOFF_S", 1)) * time.Second
	RedisMaxRetryBackoffS = time.Duration(GetEnvAsInt("REDIS_MAX_RETRY_BACKOFF_S", 1)) * time.Second
	RedisIdleTimeoutS = time.Duration(GetEnvAsInt("REDIS_IDLE_TIMEOUT_S", 300)) * time.Second
	RedisMaxConnAgeS = time.Duration(GetEnvAsInt("REDIS_MAX_CONN_AGE_S", 1800)) * time.Second

	ElasticAPMServerURL = GetEnv("ELASTIC_APM_SERVER_URL", "")
	ElasticAPMServiceName = GetEnv("ELASTIC_APM_SERVICE_NAME", "service-order")
	ElasticAPMServiceVersion = GetEnv("ELASTIC_APM_SERVICE_VERSION", "0.1.0")
	ElasticAPMEnvironment = GetEnv("ELASTIC_APM_ENVIRONMENT", "development")
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := strings.ToLower(GetEnv(name, ""))
	if valStr == "" {
		return defaultVal
	}
	return valStr == "1" || valStr == "true" || valStr == "yes"
}

func GetEnvAsInt(name string, defaultVal int) int {
	if val, err := strconv.Atoi(GetEnv(name, "")); err == nil {
		return val
	}
	return defaultVal
}
