# Используем официальный образ Golang
FROM golang:1.23.3

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum, чтобы собрать зависимости
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем исходный код в контейнер
COPY . ./

# Копируем config.env в контейнер
COPY config/config.env /app/config/config.env

# Копируем тесты в контейнер
COPY tests /app/tests

# Переходим в директорию с точкой входа в приложение
WORKDIR /app/cmd/server

# Компилируем приложение
RUN go build -o /wallet-service .

# Открываем порт 8080
EXPOSE 8080

# Запускаем скомпилированное приложение
CMD ["/wallet-service"]