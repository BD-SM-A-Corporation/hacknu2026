# Установка системы (Версия 1)

Этот раздел описывает процесс установки системы RailTelemetry v1 для серверов.

## Требования
* PHP 8.3+
* PostgreSQL 15+
* Node.js 20+

## Шаги установки

1. **Клонируйте репозиторий**
   ```bash
   git clone https://github.com/BD-SM-A-Corporation/hacknu2026.git
   cd hacknu2026
   ```

2. **Запустите Docker окружение**
   Система поставляется с Laravel Sail для быстрой контейнеризации:
   ```bash
   cp .env.example .env
   docker compose up -d
   ```

3. **Сборка фронтенда**
   ```bash
   npm install
   npm run build
   ```

[Вернуться на главную](/docs)
