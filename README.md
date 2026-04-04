# Locomotive Telemetry System

Единая система управления телеметрией локомотивов, состоящая из:

- **Svelte 5 + Inertia v3 Dashboard** (SPA)
- **Laravel 13 API Backend** (CRUD локомотивов, управление WS-эндпоинтами, аутентификация)
- **Go 1.26 Microservice** (Высоконагруженный брокер `gorilla/websocket`, воркер-пулы на горутинах, пакетная агрегация `LocomotiveState`)

## Требования

- Docker & Docker Compose
- Node.js & npm (для предварительной сборки фронтенда, если требуется)

## Быстрый старт через Docker

1. **Сборка фронтенда** (один раз локально)

   ```bash
   npm install
   npm run build
   ```

2. **Настройка окружения**
   Скопируйте `env` файл:

   ```bash
   cp .env.example .env
   # Убедитесь, что DB_HOST=postgres, DB_USERNAME=sail, Redis настроен
   ```

3. **Запуск всех сервисов (Nginx, PHP-FPM, PGSQL, Redis, Go)**

   ```bash
   docker compose up -d --build
   ```

4. **Инициализация БД (Ключи и миграции)**
   Выполнить внутри контейнера `php`:

   ```bash
   docker compose exec php php artisan key:generate
   docker compose exec php php artisan migrate:fresh --seed
   ```

   *Примечание: Go сервис автоматически мигрирует таблицу `telemetry_histories` при запуске.*

5. **Использование системы**
   - Перейдите в браузере: `http://localhost/` -> дашборд.
   - Микросервис Go запущен и слушает подключения по адресу: WebSockets: `ws://localhost/telemetry/`.
   - Зайдите в меню "WS Настройки" чтобы привязать интерфейс к Go-сервису (`ws://localhost/telemetry/`).

### Структура проекта

- `/resources/js/components/telemetry/*` — Компоненты UI (Модалки, спидометр, графики топлива).
- `/telemetry-service/cmd/server/main.go` — Входная точка микросервиса на Golang. Пул воркеров.
- `/docker/nginx/default.conf` — Маршрутизатор (Reverse Proxy).
