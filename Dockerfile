FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/subscription-service ./cmd/app

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /bin/subscription-service /bin/subscription-service

EXPOSE 8080

CMD ["/bin/subscription-service"]
