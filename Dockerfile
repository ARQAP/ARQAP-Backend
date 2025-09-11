FROM golang:1.25.1-alpine3.22

WORKDIR /app

# Installs Air
RUN go install github.com/air-verse/air@latest

# Copy dependencies first (cache)
COPY go.mod ./
RUN go mod download

# Expose port
EXPOSE 8080

# Start with Air
CMD ["air", "-c", ".air.toml"]