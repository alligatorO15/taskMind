package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger - глобальный экземпяр логгера приложения (для синглтон-паттерна )
var Logger *zap.SugaredLogger

// Init - инициалищирует логгер
// debug-режим - импользуется расщиренный формат вывода
// prod-режим - json формат для парсинга какими-нибудь внешнимим системами мониторинга
func Init(mode string) error {
	var cfg zap.Config

	if mode == "debug" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return err
	}

	Logger = zapLogger.Sugar()
	return nil
}

// Sync - сброс буферезированных записей логов (вызывем с defer)
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
