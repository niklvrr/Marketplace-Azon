FROM golang:1.24-alpine AS builder
RUN apk --no-cache add gcc g++ make
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o marketplace ./cmd/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates postgresql-client
WORKDIR /app
COPY --from=builder /app/marketplace .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/.env .env
EXPOSE 8080
CMD ["./marketplace"]
