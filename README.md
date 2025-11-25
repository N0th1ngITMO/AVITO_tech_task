# PR Review Service

Сервис для управления пользователями и pull request'ами с возможностью массовой деактивации и автоматического переназначения ревьюверов.

### Запуск с Make 

```bash
# Клонировать репозиторий
git clone https://github.com/N0th1ngITMO/AVITO_tech_task.git
cd AVITO_tech_task

# Запустить все сервисы
make run
```

### Запуск без Make

```bash
# Собрать и запустить
docker-compose up -d --build

# Остановить
docker-compose down
```

##  Доступные команды Make

###  Управление Docker
| Команда | Описание |
|---------|-------------|
| `make build` | Собрать Docker образ приложения |
| `make run` | Запустить все сервисы в фоне |
| `make stop` | Остановить все сервисы |
| `make restart` | Перезапустить все сервисы |
| `make clean` | Остановить сервисы и очистить Docker |

###  Логи и мониторинг
| Команда | Описание |
|---------|-------------|
| `make logs` | Показать логи приложения в реальном времени |
| `make db-logs` | Показать логи базы данных в реальном времени |
| `make status` | Показать статус всех сервисов |

###  Операции с базой данных
| Команда | Описание |
|---------|-------------|
| `make db-shell` | Подключиться к оболочке PostgreSQL |
| `make db-backup` | Создать резервную копию базы данных |

###  Быстрые команды
| Команда | Описание |
|---------|-------------|
| `make dev` | Запустить и следить за логами (run + logs) |
| `make rebuild` | Полная пересборка (clean + build + run) |
| `make help` | Показать все доступные команды |

## Доступ к сервисам

После запуска сервисы доступны по адресам:

- **REST API**: http://localhost:8080
- **База данных**: localhost:5432
- **Swagger документация**: http://localhost:8080/swagger/index.html

##  Настройка окружения

Перед первым запуском создайте файл `.env` в корне проекта:

```env
DB_NAME=pr_review_db
DB_USER=pr_user
DB_PASSWORD=secure_password_123

DB_HOST=localhost
DB_PORT=5432
DB_USER=pr_user
DB_PASSWORD=secure_password_123
DB_NAME=pr_review_db
DB_SSLMODE=disable
```

##  API Endpoints

### Массовая деактивация пользователей
```http
POST /users/massDeactivate
Content-Type: application/json

{
  "team_name": "backend",
  "exclude_user_ids": ["user1", "user2"]
}
```

### Создать PR
```http
POST /pullRequest/create
Content-Type: application/json

{
  "author_id": "userId",
  "pull_request_id": "pullRequestId",
  "pull_request_name": "string"
}
```

### Merge PR
```http
POST /pullRequest/merge
Content-Type: application/json

{
  "pull_request_id": "pullRequestId"
}
```

### Переназначить конкретного ревьювера
```http
POST /pullRequest/reassign
Content-Type: application/json

{
  "old_reviewer_id": "oldReviewerId",
  "pull_request_id": "pullRequestId"
}
```

### Получить статистику по PR
```http
GET /stats/prs
```

### Получить общеую стистику
```http
GET /stats/overall
```

###  Получить статистику по ревьюверам
```http
GET /stats/users
```

### Создать команду с участниками
```http
POST /team/add
Content-Type: application/json

{
  "team_name": "strongers",
  "members": [
    {
      "user_id": "u6",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "u7",
      "username": "Bob",
      "is_active": true
    },
    {
      "user_id": "u8",
      "username": "Bob",
      "is_active": true
    },
    {
      "user_id": "u9",
      "username": "Bob",
      "is_active": true
    }
  ]
}
```

###  Получить комаду с участниками
```http
GET /team/get
```

###  Получить PR, где пользователь назначен ревьювером
```http
GET /users/getReview
```

### Массовая деактивация пользователей команды
```http
POST /users/massDeactivate
Content-Type: application/json

{
  "exclude_user_ids": ["u8"],
  "team_name": "backend"
}
```

### Установить флаг активности пользователя
```http
POST /users/setIsActive
Content-Type: application/json

{
  "is_active": true,
  "user_id": "userId"
}
```

##  Архитектура

```
├── cmd/main.go          # Точка входа
├── internal/
|   ├── config/              # Config базы данных
|   ├── error/               # Ошибки
│   ├── handler/             # HTTP хендлеры
│   ├── service/             # Бизнес-логика
│   ├── repository/          # Работа с базой данных
│   ├── model/               # Модели данных
│   └── dto/                 # Data Transfer Objects
├── tests/
|   └── integration/         # Интеграционные тесты
├── docker-compose.yml       # Конфигурация Docker
├── Dockerfile              # Сборка приложения
├── Makefile               # Команды управления
└── init/                  # Скрипты инициализации БД
    └── init.sql
├── docs/                  # Сгенерированная документация Swagger
```

## База данных

Сервис использует PostgreSQL с автоматической инициализацией таблиц:

- **users** - таблица пользователей
- **pull_requests** - таблица pull request'ов
- **team** - таблица команд

## Разработка

**Примечание**: Для работы Make на Windows установите [Chocolatey](https://chocolatey.org/) и выполните:
```powershell
choco install make
```

## Проблемы возникшие во время решения

Сервис при отсутствии pr в базе данных на запрос GET /stats/overall отвечает ошибкой с кодом 500. Логика не нарушена, но код не очень хороший, не успел поправить до делайна
