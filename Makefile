# Указываем пути
SHELL := /bin/bash
MIGRATE := migrate
POSTGRES_URL := postgres://postgres:secret@localhost:5432/books?sslmode=disable
MIGRATIONS_DIR := migrations

# Установка зависимостей
.PHONY: install
install:
	@echo "🛠 Установка зависимостей..."
	go mod tidy

.PHONY: docs-gen
docs-gen:
	@echo "Генерируем доку умарова ))000)"
	swag init -g ./cmd/main.go -o ./docs

# Запуск базы данных через Docker Compose
.PHONY: start-db
start-db:
	@echo "🚀 Поднимаем инфраструктуру (PostgreSQL + MongoDB)..."
	docker-compose up -d postgres mongo

# Остановка базы
.PHONY: stop-db
stop-db:
	@echo "🛑 Останавливаем инфраструктуру..."
	docker-compose down

# Генерация миграции
.PHONY: new-migration
new-migration:
	@read -p "Введите название миграции: " NAME; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$NAME

# Применение миграций
.PHONY: migrate-up
migrate-up:
	@echo "🔄 Применяем миграции..."
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(POSTGRES_URL)" up

# Откат миграций
.PHONY: migrate-down
migrate-down:
	@echo "⏪ Откатываем миграции..."
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(POSTGRES_URL)" down 1

# Запуск бэкенда локально
.PHONY: run-backend-local
run-backend-local:
	@echo "🚀 Запускаем бэкенд локально..."
	go run main.go

# Запуск бэкенда в Docker
.PHONY: run-backend-docker
run-backend-docker:
	@echo "🐳 Запускаем бэкенд в контейнере..."
	docker-compose up --build backend

# Запуск тестов
.PHONY: test
test:
	@echo "🧪 Запуск тестов..."
	go test ./... -v

# Очистка проекта
.PHONY: clean
clean:
	@echo "🧹 Очистка проекта..."
	rm -rf migrations/*.sql
