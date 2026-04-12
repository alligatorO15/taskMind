package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config — главная структура конфигурации приложения
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	MongoDB   MongoDBConfig   `mapstructure:"mongodb"`
	RabbitMQ  RabbitMQConfig  `mapstructure:"rabbitmq"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Reminder  ReminderConfig  `mapstructure:"reminder"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
}

// ServerConfig — настройки HTTP-сервера
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// MongoDBConfig — настройки подключения к MongoDB
type MongoDBConfig struct {
	URI      string        `mapstructure:"uri"`
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// RabbitMQConfig — настройки подключения к RabbitMQ
type RabbitMQConfig struct {
	URI            string        `mapstructure:"uri"`
	Exchange       string        `mapstructure:"exchange"`
	Queue          string        `mapstructure:"queue"`
	ReconnectDelay time.Duration `mapstructure:"reconnect_delay"`
}

// JWTConfig — настройки токенов аутентификации
type JWTConfig struct {
	AccessSecret  string        `mapstructure:"access_secret"`
	RefreshSecret string        `mapstructure:"refresh_secret"`
	AccessTTL     time.Duration `mapstructure:"access_ttl"`
	RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
}

// ReminderConfig — настройки системы напоминаний
type ReminderConfig struct {
	DefaultBefore time.Duration `mapstructure:"default_before"`
	CheckInterval time.Duration `mapstructure:"check_interval"`
}

// WebSocketConfig — настройки WebSocket-соединений
type WebSocketConfig struct {
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	PongTimeout  time.Duration `mapstructure:"pong_timeout"`
	PingInterval time.Duration `mapstructure:"ping_interval"`
}

// Load — загружает конфигурацию из файла config.yaml и переменных окружения.
// Путь к файлу конфигурации задается через параметр configPath.
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("server.mode", "SERVER_MODE")
	_ = viper.BindEnv("mongodb.uri", "MONGODB_URI")
	_ = viper.BindEnv("mongodb.database", "MONGODB_DATABASE")
	_ = viper.BindEnv("rabbitmq.uri", "RABBITMQ_URI")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
