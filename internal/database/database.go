package database

import (
	"gorm-test/internal/config"
	"gorm-test/internal/domain"
	"gorm.io/driver/postgres"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Init 데이터베이스 초기화
func Init(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var logLevel logger.LogLevel
	if config.Get().Server.Mode == "debug" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	var err error
	db, err = gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return nil, err
	}

	// 자동 마이그레이션
	if err := db.AutoMigrate(&domain.Post{}, &domain.Comment{}); err != nil {
		return nil, err
	}

	log.Println("데이터베이스 연결 완료")
	return db, nil
}

// Get DB 인스턴스 반환
func Get() *gorm.DB {
	return db
}
