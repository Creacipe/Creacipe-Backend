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
	var dsn string
	useTLS := os.Getenv("DB_USE_TLS")

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	
	log.Printf("DB Config -> Host: %s, Port: %s, User: %s, Database: %s, TLS: %s\n", dbHost, dbPort, dbUser, dbName, useTLS)

	if useTLS == "true" {
	
		mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: dbHost,
		})

		
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		)
		log.Println("Connecting to database with TLS (online mode)...")
	} else {
		
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		)
		log.Println("Connecting to database without TLS (localhost mode)...")
	}

	
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             500 * time.Millisecond, 
			LogLevel:                  logger.Warn,           
			IgnoreRecordNotFoundError: true,                   
			Colorful:                  true,
		},
	)

	database, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database! Error: %v\n", err)
	}
	log.Println("Koneksi ke database berhasil.")

	DB = database
}

func GetDB() *gorm.DB {
	return DB
}