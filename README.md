# Subscription
REST API сервис для агрегации данных об онлайн подписках

## Стек
- Golang
- Postgresql
- Docker
- Swagger
- Logrus

# 🚀 Быстрый старт
### Клонирование репозитория
```bash
git clone https://github.com/ShekleinAleksey/subscriptions.git
```
```bash
cd subscriptions
```

### Запуск миграций
```bash
make migrate
```
### Генерация Swagger документации
```bash
make swag
```
### Запуск приложения
```bash
make run
```

## Конфигурация приложения
Пример файла .env
```bash
DB_USERNAME="admin"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="subscriptiondb"
DB_SSLMODE="disable"
DB_PASSWORD="root123"
```
# 📝 API ENDPOINTS
### Создание подписки
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Amediateka",
    "price": 600,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "11-2025"
  }'
```
### Получение списка подписок
```bash
curl "http://localhost:8080/api/v1/subscriptions?limit=10&offset=0"
```
### Обновление подписки
```bash
curl -X PUT http://localhost:8080/api/v1/subscriptions/a1b2c3d4-e5f6-7890-abcd-ef1234567890 \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Amediateka",
    "price": 500,
    "start_date": "12-2025",
    "end_date": "01-2026"
  }'
```
### Удаление подписки
```bash
curl -X DELETE http://localhost:8080/api/v1/subscriptions/a1b2c3d4-e5f6-7890-abcd-ef1234567890
```
### Получение суммарной стоимости
```bash
# Все подписки
curl "http://localhost:8080/api/v1/subscriptions/summary"

# По конкретному пользователю
curl "http://localhost:8080/api/v1/subscriptions/summary?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"

# По периоду и сервису
curl "http://localhost:8080/api/v1/subscriptions/summary?start_period=11-2025&end_period=12-2025&service_name=Amediateka"
```
#### Swagger документация доступна после запуска: http://localhost:8080/swagger/index.html

