FROM golang:1.25.1-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git bash
COPY go.mod go.sum go.work ./

COPY http/ ./http/
COPY collector/ ./collector/
COPY notification/ ./notification/
COPY platform/ ./platform/

RUN go work sync && go mod download all

WORKDIR /app/collector
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:3.20
WORKDIR /app

# install chromium for googledp go parser
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    ca-certificates \
    dumb-init

COPY --from=builder /app/server .

ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV PATH="/usr/bin:$PATH"

# Handles PID 1 responsibilities in Docker (signal forwarding, zombie reaping, etc.).
# Forwards SIGTERM / SIGINT properly to your app.
# Ensures child processes (like Chrome) are cleaned up correctly when the container stops.
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

EXPOSE 8080
CMD ["./server"]
