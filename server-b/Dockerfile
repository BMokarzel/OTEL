FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/server-b

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/app .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
