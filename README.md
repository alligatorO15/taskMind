# TaskMind — Умный таск-менеджер с системой напоминаний

TaskMind — полноценный веб-сервис для управления задачами с интеллектуальной системой отложенных напоминаний. Пользователи могут создавать задачи, привязывать их к проектам, назначать теги и дедлайны. Система автоматически отправляет напоминания через WebSocket и отслеживает просроченные задачи.

## Архитектура

Проект построен на принципах **Clean Architecture** с чётким разделением на слои:

```
backend/
├── cmd/server/          # Точка входа приложения
├── internal/
│   ├── config/          # Конфигурация (Viper)
│   ├── domain/
│   │   ├── models/      # Доменные модели (User, Task, Project, Notification)
│   │   └── repository/  # Интерфейсы репозиториев
│   ├── usecase/         # Бизнес-логика (Use Cases)
│   ├── delivery/
│   │   ├── http/        # HTTP-хендлеры, middleware, маршрутизация (Gin)
│   │   └── websocket/   # WebSocket-хаб для real-time уведомлений
│   ├── repository/
│   │   └── mongo/       # MongoDB-реализация репозиториев
│   ├── worker/          # Фоновые воркеры (напоминания, проверка дедлайнов)
│   └── infrastructure/  # Подключения к внешним сервисам (MongoDB, RabbitMQ, Logger)
└── pkg/
    └── apperror/        # Общие типы ошибок

frontend/
├── src/
│   ├── components/      # React-компоненты (Layout, TaskCard, TaskDialog, NotificationPanel)
│   ├── pages/           # Страницы (Login, Register, Dashboard, Tasks, Projects)
│   ├── store/           # Redux Toolkit (slices: auth, task, project, notification)
│   └── services/        # API-клиент (Axios), WebSocket-сервис
└── public/
```

## Технологический стек

### Бэкенд
| Компонент | Технология |
|-----------|-----------|
| Язык | Go 1.26 |
| Веб-фреймворк | Gin |
| База данных | MongoDB |
| Брокер сообщений | RabbitMQ (delayed message exchange) |
| WebSocket | gorilla/websocket |
| Аутентификация | JWT (access + refresh токены) |
| Конфигурация | Viper |
| Логирование | Zap |

### Фронтенд
| Компонент | Технология |
|-----------|-----------|
| Фреймворк | React 18 |
| Состояние | Redux Toolkit |
| UI | Material-UI (MUI) 5 |
| HTTP-клиент | Axios |
| WebSocket | Native WebSocket API |

### Инфраструктура
| Компонент | Технология |
|-----------|-----------|
| Контейнеризация | Docker + Docker Compose |
| Reverse Proxy | Nginx |
| CI/CD | GitHub Actions |

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.22+ (для локальной разработки)
- Node.js 20+ (для локальной разработки фронтенда)

### Запуск через Docker Compose

```bash
# Клонируем репозиторий
git clone https://github.com/username/taskmind.git
cd taskmind

# Запускаем все сервисы
docker-compose up -d

# Приложение доступно:
# - Фронтенд: http://localhost:3000
# - API: http://localhost:8080
# - RabbitMQ Management: http://localhost:15672 (guest/guest)
```

### Локальная разработка

#### Бэкенд

```bash
# Запускаем MongoDB и RabbitMQ
docker-compose up -d mongodb rabbitmq

# Переходим в директорию бэкенда
cd backend

# Устанавливаем зависимости
go mod download

# Запускаем сервер
go run cmd/server/main.go
```

#### Фронтенд

```bash
cd frontend

# Устанавливаем зависимости
npm install

# Запускаем dev-сервер
npm start
```

### Запуск тестов

```bash
cd backend
go test -v ./...
```

## API Endpoints

### Аутентификация
| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/v1/auth/register` | Регистрация |
| POST | `/api/v1/auth/login` | Авторизация |
| POST | `/api/v1/auth/refresh` | Обновление токена |
| GET | `/api/v1/profile` | Профиль пользователя |

### Задачи
| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/v1/tasks` | Создать задачу |
| GET | `/api/v1/tasks` | Список задач (с фильтрами) |
| GET | `/api/v1/tasks/:id` | Получить задачу |
| PUT | `/api/v1/tasks/:id` | Обновить задачу |
| DELETE | `/api/v1/tasks/:id` | Удалить задачу |

### Проекты
| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/v1/projects` | Создать проект |
| GET | `/api/v1/projects` | Список проектов |
| GET | `/api/v1/projects/:id` | Получить проект |
| PUT | `/api/v1/projects/:id` | Обновить проект |
| DELETE | `/api/v1/projects/:id` | Удалить проект |

### Уведомления
| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/notifications` | Список уведомлений |
| GET | `/api/v1/notifications/unread-count` | Количество непрочитанных |
| PUT | `/api/v1/notifications/:id/read` | Пометить прочитанным |
| PUT | `/api/v1/notifications/read-all` | Прочитать все |

### WebSocket
| Путь | Описание |
|------|----------|
| `GET /api/v1/ws` | WebSocket для real-time уведомлений |

## Примеры запросов

### Регистрация
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret123"}'
```

### Создание задачи с дедлайном
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "title": "Подготовить отчёт",
    "description": "Квартальный отчёт по продажам",
    "priority": "high",
    "tags": ["работа", "отчёт"],
    "deadline": "2026-03-01T18:00:00Z",
    "reminder_before": "1h"
  }'
```

### Создание проекта
```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"name":"Мой проект","description":"Описание проекта","color":"#2196F3"}'
```

## Система напоминаний

TaskMind использует двухуровневую систему напоминаний:

1. **RabbitMQ Delayed Messages** — при создании задачи с дедлайном, сообщение-напоминание помещается в очередь с задержкой (`deadline - reminder_before`). Когда время наступает, воркер обрабатывает сообщение.

2. **Deadline Checker** — фоновый воркер периодически (по умолчанию каждую минуту) сканирует базу данных на наличие просроченных задач и задач, требующих напоминания. Это обеспечивает надёжность даже при перезапуске сервиса.

Уведомления доставляются:
- Через **WebSocket** в реальном времени (если пользователь онлайн)
- Сохраняются в **MongoDB** для последующего просмотра

## Конфигурация

Настройки задаются в `backend/config.yaml`:

```yaml
server:
  port: 8080
  mode: debug        # debug | release

mongodb:
  uri: mongodb://localhost:27017
  database: taskmind

rabbitmq:
  uri: amqp://guest:guest@localhost:5672/

jwt:
  access_ttl: 15m    # Время жизни access-токена
  refresh_ttl: 168h  # Время жизни refresh-токена (7 дней)

reminder:
  default_before: 30m  # За сколько до дедлайна напоминать
  check_interval: 1m   # Интервал проверки просроченных задач
```

## Лицензия

MIT
