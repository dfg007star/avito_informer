FROM golang:1.25.1-alpine AS builder
WORKDIR /app

COPY go.mod go.sum go.work go.work.sum ./

COPY http/ ./http/
COPY collector/ ./collector/
COPY notification/ ./notification/
COPY platform/ ./platform/

RUN go work sync && go mod download all

WORKDIR /app/http
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/server .
COPY http/migrations ./migrations
COPY http/html ./html
COPY http/static ./static

EXPOSE 8081
CMD ["./server"]
