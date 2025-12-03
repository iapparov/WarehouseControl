package config

import (
	"time"
)

type AppConfig struct {
	ServerConfig   ServerConfig   `mapstructure:"server"`
	LoggerConfig   loggerConfig   `mapstructure:"logger"`
	DBConfig       dbConfig       `mapstructure:"db_config"`
	RetrysConfig   RetrysConfig   `mapstructure:"retry_strategy"`
	GinConfig      ginConfig      `mapstructure:"gin"`
	JwtConfig      JwtConfig      `mapstructure:"jwt"`
	UserConfig     UserConfig     `mapstructure:"username_config"`
	PasswordConfig PasswordConfig `mapstructure:"password_config"`
	ItemConfig     ItemConfig     `mapstructure:"item_config"`
}

type RetrysConfig struct {
	Attempts int           `mapstructure:"attempts" default:"3"`
	Delay    time.Duration `mapstructure:"delay" default:"1s"`
	Backoffs float64       `mapstructure:"backoffs" default:"2"`
}

type ginConfig struct {
	Mode string `mapstructure:"mode" default:"debug"`
}

type ServerConfig struct {
	Host string `mapstructure:"host" default:"localhost"`
	Port int    `mapstructure:"port" default:"8080"`
}

type loggerConfig struct {
	Level string `mapstructure:"level" default:"info"`
}

type postgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string
	Password string
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode" default:"disable"`
}

type dbConfig struct {
	Master          postgresConfig   `mapstructure:"postgres"`
	Slaves          []postgresConfig `mapstructure:"slaves"`
	MaxOpenConns    int              `mapstructure:"max_open_conns"`
	MaxIdleConns    int              `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration    `mapstructure:"conn_max_lifetime"`
}

type JwtConfig struct {
	JwtExpAccessToken  int `mapstructure:"jwt_exp_access_token"`
	JwtExpRefreshToken int `mapstructure:"jwt_exp_refresh_token"`
	JwtAccessSecret    string
	JwtRefreshSecret   string
}

type UserConfig struct {
	MinLength         int    `mapstructure:"min_length"`
	MaxLength         int    `mapstructure:"max_length"`
	AllowedCharacters string `mapstructure:"allowed_characters"`
	CaseInsesitive    bool   `mapstructure:"case_insensitive"`
}

type PasswordConfig struct {
	MinLength    int  `mapstructure:"min_length"`
	MaxLength    int  `mapstructure:"max_length"`
	RequireUpper bool `mapstructure:"require_upper"`
	RequireLower bool `mapstructure:"require_lower"`
	RequireDigit bool `mapstructure:"require_digit"`
}

type ItemConfig struct {
	NameMinLength int `mapstructure:"name_min_length"`
	NameMaxLegth  int `mapstructure:"name_max_length"`
}
