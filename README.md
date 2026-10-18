# ИИ-бот учёта семейных расходов

Сервис для удобного учёта семейных расходов, состоящий из умного Telegram-бота с ИИ-агентом и веб-панели управления для главы семьи.

---

## 🚀 Возможности

### 1. Telegram-бот с ИИ-агентом
- **Естественный язык**: общение без кнопок и команд — бот автоматически определяет сумму, категорию, дату и комментарий.
- **Поддержка полного CRUD**:
  - ➕ **Добавление**: `кофе 350`, `такси 1200 вчера`, `продукты в супермаркете 2450`
  - 🔍 **Просмотр и расчёт сумм**: `сколько потратил на еду в этом месяце`, `покажи траты за неделю`
  - ✏️ **Изменение**: `последнее было не 350, а 450`
  - 🗑️ **Удаление**: `удали такси`, `отмени последнюю трату`
- **Изоляция данных**: каждому члену семьи доступны исключительно его личные траты.
- **Гибкость LLM**: поддержка любого OpenAI-совместимого провайдера (OpenAI, DeepSeek, Groq, Ollama, Mistral и др.), настраиваемого через переменные окружения.

### 2. Веб-панель для главы семьи
- Единое окно аналитики и управления расходами всех членов семьи (`http://localhost:3000`).
- **Сводные метрики**: общая сумма трат, количество записей, число активных участников семьи.
- **Гибкая фильтрация**:
  - по члену семьи (Telegram ID);
  - по категории расхода;
  - по произвольному периоду дат (с быстрыми пресетами: «Сегодня», «7 дней», «Этот месяц», «За всё время»).
- **Полное управление**: просмотр списка, добавление, редактирование и удаление трат.
- **Современный UI**: адаптивный интерфейс в тёмной теме со стеклянным дизайном (Glassmorphism), цветовой маркировкой категорий и анимациями.
- **Без авторизации**: доступна сразу при открытии (согласно спецификации).

---

## 🛠 Стек технологий

| Компонент | Технологии |
|---|---|
| **Бэкенд** | Go 1.23+, Gin Web Framework, GORM |
| **База данных** | PostgreSQL 16 |
| **Фронтенд** | Vue 3 (Composition API), Vite, Vanilla CSS |
| **Веб-сервер** | Nginx (раздача статики SPA + reverse proxy к API) |
| **ИИ / LLM** | OpenAI-совместимый Chat Completions API (JSON Mode via System Prompt) |
| **Telegram API** | `go-telegram-bot-api/v5` |
| **Контейнеризация** | Docker, Docker Compose (multi-stage builds) |

---

## 📂 Структура проекта

```
ai-bot-expenses/
├── docker-compose.yml        # Конфигурация запуска всех сервисов
├── .env.example              # Пример переменных окружения
├── .gitignore
├── README.md
│
├── backend/
│   ├── Dockerfile            # Multi-stage сборка Go-сервера
│   ├── go.mod, go.sum
│   ├── cmd/server/main.go    # Точка входа: параллельный запуск HTTP API и Telegram-бота
│   └── internal/
│       ├── config/           # Загрузка и валидация конфигурации
│       ├── models/           # Модели данных и валидаторы категорий
│       ├── db/               # Подключение к PostgreSQL и AutoMigrate
│       ├── repository/       # Слой работы с БД (GORM)
│       ├── service/          # Бизнес-логика учёта трат
│       ├── handler/          # REST API эндпоинты (Gin) и CORS
│       ├── bot/              # Telegram-бот и обработчик сообщений
│       └── llm/              # HTTP-клиент к LLM и системный промпт
│
└── frontend/
    ├── Dockerfile            # Сборка Vue 3 и отдача через Nginx
    ├── nginx.conf            # Конфигурация Nginx (SPA + прокси /api/)
    ├── package.json, vite.config.js
    ├── index.html
    └── src/
        ├── App.vue           # Главный компонент веб-панели
        ├── api/expenses.js   # Клиент к REST API бэкенда
        ├── assets/main.css   # Дизайн-система (CSS-переменные, glassmorphism)
        └── components/
            ├── AppHeader.vue      # Шапка со статистикой
            ├── ExpenseFilters.vue # Панель фильтров
            ├── ExpenseTable.vue   # Таблица трат с сортировкой
            ├── ExpenseForm.vue    # Модальное окно создания/редактирования
            └── ConfirmDialog.vue  # Подтверждение удаления
```

---

## 📋 Модель данных и категории

Каждая запись расхода содержит:
- `id`: уникальный идентификатор
- `user_id`: Telegram ID владельца траты
- `amount`: сумма расхода
- `category`: одна из 7 стандартных категорий (`еда`, `транспорт`, `жильё`, `развлечения`, `здоровье`, `покупки`, `другое`)
- `spent_at`: дата траты (`YYYY-MM-DD`)
- `comment`: комментарий / наименование траты
- `created_at`, `updated_at`: системные метки времени

---

## ⚙️ Настройка и запуск

### 1. Подготовка конфигурации

Скопируйте пример файла конфигурации:
```bash
cp .env.example .env
```

Заполните переменные в `.env`:
```env
# Telegram-бот (получить у @BotFather)
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrSTUvwxYZ

# LLM-провайдер (OpenAI, DeepSeek, Groq, Ollama и др.)
LLM_BASE_URL=https://api.openai.com/v1
LLM_API_KEY=sk-your-llm-api-key-here
LLM_MODEL=gpt-4o-mini

# Настройки БД (можно оставить по умолчанию)
DB_HOST=postgres
DB_PORT=5432
DB_USER=expenses
DB_PASSWORD=expenses
DB_NAME=expenses
```

> **Пример для DeepSeek:**
> `LLM_BASE_URL=https://api.deepseek.com/v1`
> `LLM_MODEL=deepseek-chat`
>
> **Пример для локальной Ollama:**
> `LLM_BASE_URL=http://host.docker.internal:11434/v1`
> `LLM_API_KEY=ollama`
> `LLM_MODEL=llama3.1`

### 2. Запуск через Docker Compose

Запустите проект одной командой:
```bash
docker compose up --build -d
```

После старта будут доступны:
- 🌐 **Веб-панель**: [http://localhost:3000](http://localhost:3000)
- 🔌 **REST API**: [http://localhost:8080/api/expenses](http://localhost:8080/api/expenses)
- 🤖 **Telegram-бот**: готов к приёму сообщений в Telegram

Для просмотра логов:
```bash
docker compose logs -f
```

Для остановки:
```bash
docker compose down
```

---

## 🔌 REST API

| Метод | Эндпоинт | Описание | Параметры запроса |
|---|---|---|---|
| `GET` | `/api/expenses` | Список трат с фильтрацией | `user_id`, `category`, `date_from`, `date_to` |
| `POST` | `/api/expenses` | Добавить новую трату | Тело JSON: `user_id`, `amount`, `category`, `spent_at`, `comment` |
| `GET` | `/api/expenses/:id` | Получить трату по ID | — |
| `PUT` | `/api/expenses/:id` | Редактировать трату | Тело JSON (частичное обновление полей) |
| `DELETE` | `/api/expenses/:id` | Удалить трату | — |
| `GET` | `/api/users` | Список ID всех участников | — |
| `GET` | `/health` | Проверка работоспособности | — |
