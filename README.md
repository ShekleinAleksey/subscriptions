# Subscription
REST API сервис для агрегации данных об онлайн подписках

## Стек
- Golang
- Postgresql
- Docker

# 🚀 Быстрый старт
### Запуск с Docker
```bash
git clone https://github.com/ShekleinAleksey/subscriptions.git
```
```bash
cd subscriptions
```
```bash
docker-compose up
```

## Конфигурация приложения
Пример файла config.env
```bash
DB_USERNAME="admin"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="subscriptiondb"
DB_SSLMODE="disable"
DB_PASSWORD="root123"
PG_URL="postgres://admin:root123@localhost:5432/subscriptiondb"
```
## API ENDPOINTS

### Тестирование
```bash
go test ./... -v
