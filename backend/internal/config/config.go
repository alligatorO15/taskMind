package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config - главная структура конфига приложения
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	MongoDB   MongoDBConfig   `mapstructure:"monodb"`
	RabbitMQ  RabbitMQConfig  `mapstructure:"rabbitmq"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Reminder  ReminderConfig  `mapstructure:"reminder"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
}

// настройки http-сервера
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// настройки подключений к mongodb
type MongoDBConfig struct {
	URI      string        `mapstructure:"uri"`
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// настройки подключения к rabbitmq
type RabbitMQConfig struct {
	URI            string        `mapstructure:"uri"`
	Exchange       string        `mapstructure:"exchange"`
	Queue          string        `mapstructure:"queue"`
	ReconnectDelay time.Duration `mapstructure:"reconnect_delay"`
}

// настройки токенов аутентификации
type JWTConfig struct {
	AccessSecret  string        `mapstructure:"access_secret"`
	RefreshSecret string        `mapstructure:"refresh_secret"`
	AccessTTL     time.Duration `mapstructure:"access_ttl"`
	RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
}

// настройки системы напоминаний
type ReminderConfig struct {
	DefaultBefore time.Duration `mapstructure:"default_before"`
	CheckInterval time.Duration `mapstructure:"check_interval"`
}

// настройки web-socket соединений
type WebSocketConfig struct {
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	PingInterval time.Duration `mapstructure:"ping_interval"`
	PongTimeout  time.Duration `mapstructure:"pong_timeout"`
}

// Load - загружает конфиг из config.yml и переменных окружения
// путь к файлу задается через параметр configPath
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
