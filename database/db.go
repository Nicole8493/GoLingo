package db

import (
	"fmt"
	"github.com/Nicole8493/GoLingo/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func New(config config.Config) (*gorm.DB, error) {
	connectString := ""
	var db *gorm.DB
	var err error

	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	if config.DB.Debug {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
			},
		)
		dbConfig.Logger = newLogger
	}

	connectString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
	)
	if config.DB.SSLMode != "" {
		connectString += " sslmode=" + config.DB.SSLMode
	}
	db, err = gorm.Open(postgres.Open(connectString), dbConfig)

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		Article{},
		Translation{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
