# 📈 Currency Tracker Service

English | [Русский](README.ru.md)

An asynchronous REST API microservice built in Go for automated cryptocurrency tracking, thread-safe in-memory caching, and historical data persistence with JWT-based authentication.

---

## 🚀 Key Features & Architecture

* **Concurrent Background Worker (Fan-Out / Fan-In):** Parallel polling of external CoinGecko endpoints using goroutines, `sync.WaitGroup`, and buffered channels without blocking the HTTP server.
* **Layered Architecture:** Strict separation of concerns across API handlers, domain logic, in-memory caching, and SQLite persistence.
* **Thread-Safe In-Memory Cache:** Custom concurrent cache implementation powered by `sync.RWMutex` for low-latency price lookups.
* **Data Persistence (SQLite):** Historical quotes and user credential storage using standard `database/sql` with parameterized queries to prevent SQL injections.
* **Auth & Security:** User registration and login flow, password hashing via `bcrypt` with salt, session validation via `JWT` (supporting both `Authorization: Bearer` headers and browser-based `HttpOnly` cookies).
* **Graceful Shutdown:** Safe termination of background workers and HTTP listeners on system interrupt signals (`SIGINT`, `SIGTERM`) via `context.Context`.
* **Testing:** Unit-tested HTTP client and cache layers using `net/http/httptest` to mock remote network responses.

---

## 🛠 Tech Stack

* **Language:** Go (Golang 1.22+)
* **Database:** SQLite (`modernc.org/sqlite` — pure Go without CGO)
* **Authentication:** JWT (`golang-jwt/jwt/v5`), Bcrypt (`golang.org/x/crypto/bcrypt`)
* **Configuration:** YAML (`gopkg.in/yaml.v3`)
* **Testing:** `testing`, `net/http/httptest`

---

## 📁 Project Structure

```text
.
├── config.yaml               # Runtime configuration (ports, intervals, DB path)
├── main.go                   # Dependency injection and application bootstrap
├── internal/
│   ├── api/                  # HTTP handlers and authentication middleware
│   ├── auth/                 # JWT token generation and bcrypt password hashing
│   ├── cache/                # Thread-safe in-memory cache
│   ├── client/               # External CoinGecko API client + unit tests
│   ├── config.go             # YAML configuration parsing and validation
│   ├── storage/              # SQLite persistence (schema initialization, queries)
│   └── worker/               # Background concurrent polling worker

# 📈 Currency Tracker Service

Асинхронный REST API микросервис на Go для регулярного мониторинга, кеширования и сохранения котировок криптовалют с защищенным доступом по JWT.

---

## 🚀 Ключевые возможности и архитектура

* **Многопоточный фоновый воркер (Fan-Out / Fan-In):** параллельный сбор котировок из внешнего API (CoinGecko) с использованием горутин, `sync.WaitGroup` и буферизированных каналов без блокировки веб-сервера.
* **Слоистая архитектура (Clean Architecture):** строгое разделение на слои API-хэндлеров, бизнес-логики, in-memory кеша и хранилища (SQLite).
* **Потокобезопасный RAM-кеш:** собственная реализация кеша с конкурентным доступом на базе `sync.RWMutex` для мгновенной отдачи последних цен.
* **Персистентность (SQLite):** сохранение исторических срезов цен и пользователей через стандартный пакет `database/sql` с параметризованными запросами для защиты от SQL-инъекций.
* **Безопасность (Auth & Security):** регистрация и аутентификация пользователей, хеширование паролей через `bcrypt` с солью, валидация сессий через `JWT` (поддержка заголовка `Authorization: Bearer` и браузерных `HttpOnly` Cookie).
* **Graceful Shutdown:** безопасная остановка сервера и фоновых горутин с перехватом системных сигналов (`SIGINT`, `SIGTERM`) через `context.Context`.
* **Тестирование:** покрытие юнит-тестами HTTP-клиента и кеша с мокированием сетевых ответов через `net/http/httptest`.

---

## 🛠 Стек технологий

* **Язык:** Go (Golang 1.22+)
* **База данных:** SQLite (`modernc.org/sqlite` — чистый Go без CGO)
* **Аутентификация:** JWT (`golang-jwt/jwt/v5`), Bcrypt (`golang.org/x/crypto/bcrypt`)
* **Конфигурация:** YAML (`gopkg.in/yaml.v3`)
* **Тестирование:** `testing`, `net/http/httptest`

---

## 📁 Структура проекта

```text
.
├── config.yaml               # Файл конфигурации (порты, интервалы, путь к БД)
├── main.go                   # Инициализация зависимостей и запуск сервера
├── internal/
│   ├── api/                  # HTTP-хэндлеры и middleware авторизации
│   ├── auth/                 # Генерация токенов JWT и хеширование bcrypt
│   ├── cache/                # Потокобезопасный in-memory кеш
│   ├── client/               # Клиент для работы с CoinGecko API + Unit-тесты
│   ├── config.go             # Парсинг и валидация конфигурации
│   ├── storage/              # Слой базы данных (SQLite: схемы, миграции, CRUD)
│   └── worker/               # Конкурентный фоновый воркер опроса котировок