FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /api ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
