package config

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Register TLS config for TiDB Cloud
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: os.Getenv("DB_HOST"),
	})

	// Build DSN with TLS
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Configure GORM logger for production (Warn level - only slow queries and errors)
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             500 * time.Millisecond, // Slow query threshold
			LogLevel:                  logger.Warn,            // Only log warnings and errors
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound
			Colorful:                  true,
		},
	)

	database, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic("Gagal terhubung ke database!")
	}
	log.Println("Koneksi ke database berhasil.")

	DB = database
}

func GetDB() *gorm.DB {
	return DB
}