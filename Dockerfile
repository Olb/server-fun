# build stage
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache curl build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o blog-api cmd/main.go

# migrate binary
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate

# final stage
FROM alpine:latest
WORKDIR /root
RUN apk add --no-cache ca-certificates libc6-compat

COPY --from=builder /app/blog-api .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY migrations /migrations

EXPOSE 8080

CMD ["sh", "-c", "migrate -path=/migrations -database=${POSTGRESQL_URL} up && ./blog-api"]
