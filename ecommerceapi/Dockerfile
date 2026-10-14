# syntax=docker/dockerfile:1

# 1. Aşama: Build
FROM golang:1.25.8-alpine AS builder

WORKDIR /app

# Bağımlılıkları indir (Modül önbelleğini bağlayarak)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Kaynak kodları kopyala
COPY . .

# Derleme: -a ve -installsuffix kaldırıldı, BuildKit cache mount eklendi
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o main .

# 2. Aşama: Production
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 7116

CMD ["./main"]