package config

import (
	"fmt"
	"os"

	wbfconfig "github.com/wb-go/wbf/config"
)

func NewAppConfig() (*AppConfig, error) {
	envFilePath := "./.env"
	appConfigFilePath := "./config/local.yaml"

	cfg := wbfconfig.New()

	// Загрузка .env файлов
	if err := cfg.LoadEnvFiles(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to load env files: %w", err)
	}

	// Включение поддержки переменных окружения
	cfg.EnableEnv("")

	if err := cfg.LoadConfigFiles(appConfigFilePath); err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}

	var appCfg AppConfig
	if err := cfg.Unmarshal(&appCfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	appCfg.DBConfig.Master.DBName = os.Getenv("POSTGRES_DB")
	appCfg.DBConfig.Master.User = os.Getenv("POSTGRES_USER")
	appCfg.DBConfig.Master.Password = os.Getenv("POSTGRES_PASSWORD")

	appCfg.JwtConfig.JwtAccessSecret = os.Getenv("JWT_ACCESS_SECRET")
	appCfg.JwtConfig.JwtRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")
	return &appCfg, nil
}
