FROM golang:1.25.1-alpine3.22

WORKDIR /app

# Instalar Air
RUN go install github.com/air-verse/air@latest

# Copiar dependencias primero (cache)
COPY go.mod ./
RUN go mod download

# Exponer puerto
EXPOSE 8080

# Arrancar con Air
CMD ["air", "-c", ".air.toml"]
