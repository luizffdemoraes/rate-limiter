# Etapa de build
FROM golang:1.22.3 AS builder

WORKDIR /app

# Cache de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Imagem final mínima
FROM alpine:3.19

WORKDIR /app

# Certificados raiz (caso precise fazer chamadas HTTP externas)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /app/server

# Porta conforme README
EXPOSE 8080

# Usa variáveis de ambiente do container
ENV HTTP_PORT=8080

ENTRYPOINT ["/app/server"]