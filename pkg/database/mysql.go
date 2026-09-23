package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"survey-backend/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB opens the pool and verifies it before returning, so a bad
// configuration fails at boot rather than on the first request.
func ConnectDB() error {
	host := config.Get("DB_HOST", "127.0.0.1")
	user := config.Get("DB_USER", "root")
	password := config.Get("DB_PASSWORD", "")
	dbName := config.Get("DB_NAME", "")
	port := config.Get("DB_PORT", "3306")

	// MySQL DSN format:
	// username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	params := "charset=utf8mb4&parseTime=True&loc=Local"

	// Hosted MySQL (TiDB Cloud, Aiven, PlanetScale-likes) requires TLS; a local
	// server usually does not offer it. DB_PARAMS overrides the lot for
	// providers that need something more specific.
	if config.Bool("DB_TLS", false) {
		params += "&tls=true"
	}
	if override := config.Get("DB_PARAMS", ""); override != "" {
		params = override
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, dbName, params)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger(),
		// TiDB and other MySQL-compatible engines handle foreign keys
		// differently from MySQL, which can make AutoMigrate fail on the
		// ON DELETE CASCADE constraints in the models. The cascades are
		// belt-and-braces anyway: every child row is deleted explicitly in Go.
		DisableForeignKeyConstraintWhenMigrating: config.Bool("DB_DISABLE_FK", false),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to reach the connection pool: %w", err)
	}

	// Managed MySQL plans cap concurrent connections, often in the low tens.
	// Leaving these unset lets the pool grow until the provider refuses.
	sqlDB.SetMaxOpenConns(config.Int("DB_MAX_OPEN_CONNS", 10))
	sqlDB.SetMaxIdleConns(config.Int("DB_MAX_IDLE_CONNS", 5))
	sqlDB.SetConnMaxLifetime(config.Duration("DB_CONN_MAX_LIFETIME", 30*time.Minute))
	sqlDB.SetConnMaxIdleTime(config.Duration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), config.Duration("DB_CONNECT_TIMEOUT", 15*time.Second))
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database did not respond: %w", err)
	}

	DB = db
	fmt.Println("Successfully connected to MySQL database")

	return nil
}

// newLogger builds GORM's logger with "record not found" demoted. That error is
// ordinary control flow here - checking an email is free, looking for the seeded
// admin, a failed login - and GORM logs it at error level by default, which
// would put a stack of noise in the platform's log stream on every signup.
func newLogger() gormLogger.Interface {
	return gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel(),
			IgnoreRecordNotFoundError: true,
			Colorful:                  config.Get("APP_ENV", "development") != "production",
		},
	)
}

// logLevel keeps query logging quiet in production, where GORM's default
// slow-query warnings would flood the platform's log stream.
func logLevel() gormLogger.LogLevel {
	switch config.Get("DB_LOG_LEVEL", "") {
	case "silent":
		return gormLogger.Silent
	case "error":
		return gormLogger.Error
	case "warn":
		return gormLogger.Warn
	case "info":
		return gormLogger.Info
	}

	if config.Get("APP_ENV", "development") == "production" {
		return gormLogger.Error
	}

	return gormLogger.Warn
}

// Close releases the pool, used on shutdown.
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
