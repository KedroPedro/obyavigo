package config

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Environment Environment `toml:"environment"`
	Server      Server      `toml:"server"`
	Logger      Logger      `toml:"logger"`
	Database    Database    `toml:"database"`
	EMail       EMail       `toml:"mail"`
	Secure      Secure      `toml:"secure"`
}

type Environment struct {
	Env          string `toml:"env"`
	TemplatesDir string `toml:"templatesDir"`
	StaticDir    string `toml:"staticDir"`
}

type Server struct {
	Addr        string        `toml:"addr"`
	Timeout     time.Duration `toml:"timeout"`
	IdleTimeout time.Duration `toml:"idle_timeout"`
}

type Logger struct {
	ErrLogPath  string `toml:"error_log_path"`
	InfoLogPath string `toml:"info_log_path"`
}

type Database struct {
	Postgres Postgres `toml:"postgres"`
	Mongo    Mongo    `toml:"mongo"`
	Queries  string   `toml:"queriesFile"`
}

type Postgres struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

type Mongo struct {
	Host     string        `toml:"host"`
	Port     int           `toml:"port"`
	User     string        `toml:"user"`
	Password string        `toml:"password"`
	DBName   string        `toml:"dbname"`
	Timeout  time.Duration `toml:"timeout"`
}

type EMail struct {
	Username   string `toml:"username"`
	Password   string `toml:"password"`
	Port       int    `toml:"tlsport"`
	SMTPServer string `toml:"smtpserver"`
}

type Secure struct {
	SecretKey string `toml:"secretkey"`
}

const fn = "internal.config.New"

func New() *Config {
	cfg := new(Config)
	if _, err := toml.DecodeFile("./configs/config.toml", cfg); err != nil {
		log.Fatal("config file parsing error", fn)
	}
	return cfg
}
