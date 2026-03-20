package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Nicole8493/GoLingo/config"
	"github.com/Nicole8493/GoLingo/controller"
	db "github.com/Nicole8493/GoLingo/database"
	"github.com/Nicole8493/GoLingo/usecase"
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	// Parse command line flags получили путь до файла конфиг
	configPath := flag.String("config", "./cfg/config.yml", "path to config file (optional)")
	flag.Parse()

	// Initialize logger
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// Parse config
	cfg, err := loadConfig(*configPath)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Initialize database
	db, err := db.New(cfg)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}

	useCase := usecase.New(db, cfg.PrivateKey)

	controller, err := controller.New(cfg, useCase, db, cfg.PrivateKey)
	if err != nil {
		logger.Fatal("failed to initialize controller", zap.Error(err))
	}

	err = controller.Run(context.Background())
	if err != nil {
		logger.Fatal("failed to initialize controller", zap.Error(err))
	}

}

// loadConfig loads configuration from YAML file and environment variables
func loadConfig(configPath string) (config.Config, error) {
	var cfg config.Config

	// Load config from YAML file with environment variable override support
	err := configor.Load(&cfg, configPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to load config from %s: %w", configPath, err)
	}

	return cfg, nil
}
