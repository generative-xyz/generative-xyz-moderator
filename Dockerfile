FROM golang:1.18 as builder

WORKDIR /app

COPY ./ ./

RUN set -eux; \
    go mod download

RUN go build -o renderinghub-server cmd/main.go

## Today ubuntu is using minimalized image by default, using ubuntu for better compatible than alpine
FROM ubuntu:20.04
RUN apt-get update && apt-get install -y ca-certificates wget

WORKDIR /app

EXPOSE 10000
EXPOSE 8000

COPY --from=builder /app/renderinghub-server ./
COPY --from=builder /app/.env.production ./.env

RUN chmod +x /app/renderinghub-server
CMD ["./renderinghub-server", "app"]
