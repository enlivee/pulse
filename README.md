# Pulse

Сервис мониторинга доступности (uptime monitoring), написанный на Go.

## Описание

Pulse периодически проверяет доступность сервисов и фиксирует их состояние. Проект написан на Go с минимумом внешних зависимостей: единственная сторонняя библиотека — [`google/uuid`](https://github.com/google/uuid).

## Требования

- Go 1.27+

## Установка

```bash
git clone https://github.com/enlivee/pulse.git
cd pulse
go mod download
```

## Запуск

```bash
go run ./cmd/pulse
```

Сборка бинарника:

```bash
go build -o pulse ./cmd/pulse
./pulse
```

## Структура проекта

```
.
├── cmd/pulse/   # точка входа
├── internal/    # внутренняя логика сервиса
├── go.mod
└── go.sum
```
