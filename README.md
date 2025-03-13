# no-fines
test task

для запуска проекта неоюходимо запустить контейнер с приложением и бд
`docker-compose up --build` - применятся миграции и запустится приложение

1. HTTP API
Получение курса обмена валют
URL : /exchange/rate
Метод : POST
Тело запроса :
```json
{
  "base_currency": "RUB",
  "quote_currency": "USD"
}
```

Пример ответа:
```json
{
  "rate": 0.0123
}
```

Пример курла:
```
curl -X POST http://localhost:8080/exchange/rate \
  -H "Content-Type: application/json" \
  -d '{"base_currency": "RUB", "quote_currency": "USD"}'
```

Аналогично для gRPC

2. gRPC API
Получение курса обмена валют
Тело запроса :
```json
{
  "base_currency": "RUB",
  "quote_currency": "USD"
}
```

Пример ответа:
```json
{
  "rate": 0.0123
}
```

Пример курла:
```
grpcurl -plaintext \
  -d '{"base_currency": "RUB", "quote_currency": "USD"}' \
  localhost:50051 nofines.ExchangeService/GetExchangeRate
```
