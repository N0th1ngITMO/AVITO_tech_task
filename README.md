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
| `make run-dev` | Запустить сервисы с выводом логов |
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

- **REST API**: http://localhost:8081
- **База данных**: localhost:5432
- **Swagger документация**: http://localhost:8081/swagger/index.html

##  Настройка окружения

Перед первым запуском создайте файл `.env` в корне проекта:

```env
DB_NAME=pr_review_db
DB_USER=postgres
DB_PASSWORD=your_secure_password
SERVER_PORT=8080
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

### Получить всех пользователей
```http
GET /users
```

### Создать PR
```http
POST /pr
Content-Type: application/json

{
  "author_id": "user123",
  "title": "New feature implementation"
}
```

##  Архитектура

```
├── cmd/main.go          # Точка входа
├── internal/
│   ├── handler/             # HTTP хендлеры
│   ├── service/             # Бизнес-логика
│   ├── repository/          # Работа с базой данных
│   ├── model/               # Модели данных
│   └── dto/                 # Data Transfer Objects
├── docker-compose.yml       # Конфигурация Docker
├── Dockerfile              # Сборка приложения
├── Makefile               # Команды управления
└── init/                  # Скрипты инициализации БД
    └── init.sql
```

## База данных

Сервис использует PostgreSQL с автоматической инициализацией таблиц:

- **users** - таблица пользователей
- **pull_requests** - таблица pull request'ов

## Разработка

### Локальная разработка

```bash
make dev          # Запуск в режиме разработки
make test         # Запуск тестов
make restart      # Перезапуск после изменений
```

**Примечание**: Для работы Make на Windows установите [Chocolatey](https://chocolatey.org/) и выполните:
```powershell
choco install make
```
