package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    int
	GinMode string
	DB      DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Loc      string

	MaxOpenConns       int
	MaxIdleConns       int
	ConnMaxLifetimeMin int
}

func Load() Config {
	// Local development convenience: if `.env` exists, load it.
	// In production you typically rely on real environment variables.
	_ = godotenv.Load()

	port := atoiOrDefault(getenv("PORT", "8080"), 8080)
	ginMode := getenv("GIN_MODE", "debug")
	dbPort := atoiOrDefault(getenv("DB_PORT", "3306"), 3306)

	maxOpen := atoiOrDefault(getenv("DB_MAX_OPEN_CONNS", "25"), 25)
	maxIdle := atoiOrDefault(getenv("DB_MAX_IDLE_CONNS", "25"), 25)
	lifetimeMin := atoiOrDefault(getenv("DB_CONN_MAX_LIFETIME_MINUTES", "60"), 60)

	return Config{
		Port:    port,
		GinMode: ginMode,
		DB: DBConfig{
			Host:               getenv("DB_HOST", "127.0.0.1"),
			Port:               dbPort,
			User:               getenv("DB_USER", "root"),
			Password:           getenv("DB_PASSWORD", ""),
			Name:               getenv("DB_NAME", "ecommerce"),
			Loc:                getenv("DB_LOC", "UTC"),
			MaxOpenConns:       maxOpen,
			MaxIdleConns:       maxIdle,
			ConnMaxLifetimeMin: lifetimeMin,
		},
	}
}

func (d DBConfig) DSN() string {
	// go-sql-driver/mysql DSN format:
	//   user:password@tcp(host:port)/dbname?param=value
	// We set parseTime=true so DATETIME/TIMESTAMP scan into time.Time.
	// interpolateParams=true is helpful for correct logging/debugging with placeholders.
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=%s&interpolateParams=true",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.Loc,
	)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func atoiOrDefault(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
