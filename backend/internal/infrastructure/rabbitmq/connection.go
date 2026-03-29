package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/alligatorO15/taskMind/backend/internal/config"
	"github.com/alligatorO15/taskMind/backend/internal/domain/models"
	"github.com/alligatorO15/taskMind/backend/internal/infrastructure/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection — обёртка над подключением к RabbitMQ с поддержкой переподключения.
// Предоставляет методы для публикации отложенных сообщений и потребления из очереди.
// Использует раздельные каналы для publish и consume, т.к. amqp.Channel не потокобезопасен.
type Connection struct {
	conn      *amqp.Connection
	publishCh *amqp.Channel
	consumeCh *amqp.Channel // не нужен мьютекс,так как однократно при инициализации
	cfg       config.RabbitMQConfig
	publishMu sync.Mutex
	closed    bool
}

// NewConnection — устанавливает подключение к RabbitMQ,
// создает exchange с поддержкой отложенных сообщений и привязывает к нему очередь.
func NewConnection(cfg config.RabbitMQConfig) (*Connection, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к RabbitMQ: %w", err)
	}

	publishCh, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("не удалось открыть publish-канал RabbitMQ: %w", err)
	}

	consumeCh, err := conn.Channel()
	if err != nil {
		publishCh.Close()
		conn.Close()
		return nil, fmt.Errorf("не удалось открыть consume-канал RabbitMQ: %w", err)
	}

	// Объявляем exchange с типом x-delayed-message для отложенной доставки.
	// Требуется плагин rabbitmq_delayed_message_exchange.
	err = publishCh.ExchangeDeclare(
		cfg.Exchange,
		"x-delayed-message",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		// Если плагин не установлен, используем стандартный direct exchange
		logger.Logger.Warnf("Плагин delayed_message_exchange не найден, используем direct exchange: %v", err)
		err = publishCh.ExchangeDeclare(
			cfg.Exchange,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			publishCh.Close()
			consumeCh.Close()
			conn.Close()
			return nil, fmt.Errorf("не удалось объявить exchange: %w", err)
		}
	}

	// Объявляем очередь для напоминаний
	_, err = publishCh.QueueDeclare(
		cfg.Queue,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		publishCh.Close()
		consumeCh.Close()
		conn.Close()
		return nil, fmt.Errorf("не удалось объявить очередь: %w", err)
	}

	// Привязываем очередь к exchange
	err = publishCh.QueueBind(
		cfg.Queue,
		cfg.Queue, // routing key = имя очереди
		cfg.Exchange,
		false,
		nil,
	)
	if err != nil {
		publishCh.Close()
		consumeCh.Close()
		conn.Close()
		return nil, fmt.Errorf("не удалось привязать очередь к exchange: %w", err)
	}

	logger.Logger.Info("Подключение к RabbitMQ установлено")

	return &Connection{
		conn:      conn,
		publishCh: publishCh,
		consumeCh: consumeCh,
		cfg:       cfg,
	}, nil
}

// PublishReminder — отправляет отложенное сообщение-напоминание в очередь.
// Параметр delay определяет задержку перед доставкой сообщения получателю.
func (c *Connection) PublishReminder(ctx context.Context, msg models.ReminderMessage, delay time.Duration) error {
	c.publishMu.Lock()
	defer c.publishMu.Unlock()

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать сообщение: %w", err)
	}

	// Задержка передается через заголовок x-delay в миллисекундах
	headers := amqp.Table{}
	if delay > 0 {
		headers["x-delay"] = int64(delay / time.Millisecond)
	}

	err = c.publishCh.PublishWithContext(
		ctx,
		c.cfg.Exchange,
		c.cfg.Queue,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("не удалось опубликовать сообщение: %w", err)
	}

	logger.Logger.Infof("Напоминание отправлено в очередь с задержкой %v для задачи %s", delay, msg.TaskID)
	return nil
}

// Consume — запускает потребление сообщений из очереди напоминаний.
// Возвращает канал для чтения входящих сообщений.
func (c *Connection) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := c.consumeCh.Consume(
		c.cfg.Queue,
		"",    // consumer tag
		false, // auto-ack (подтверждаем вручную после обработки)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось подписаться на очередь: %w", err)
	}

	return msgs, nil
}

// Close — закрывает каналы и подключение к RabbitMQ
func (c *Connection) Close() error {
	c.publishMu.Lock()
	defer c.publishMu.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true

	if c.publishCh != nil {
		c.publishCh.Close()
	}
	if c.consumeCh != nil {
		c.consumeCh.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
