# Используем официальный образ Golang как базовый
FROM golang:1.21

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum в контейнер
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы в контейнер
COPY . .

# Запускаем приложение
CMD ["go", "run", "cmd/main.go", "-c=config.json"]
