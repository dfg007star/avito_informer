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

# chromium for the chromedp parser; xvfb + x11vnc + fluxbox so Chromium runs
# non-headless (PerimeterX detects headless strongly) and the operator can VNC
# in once to log into Avito and seed the profile.
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    font-noto-cjk \
    ca-certificates \
    dumb-init \
    dbus \
    dbus-x11 \
    xvfb \
    x11vnc \
    fluxbox \
    bash

COPY --from=builder /app/server .
COPY collector/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV CHROME_PATH=/usr/bin/chromium-browser
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV PATH="/usr/bin:$PATH"
ENV DISPLAY=:99
# Set VNC_ENABLED=true (and ideally VNC_PASSWORD) to expose the desktop for the
# one-time profile login. Leave unset for normal running.
ENV VNC_ENABLED=false

# dumb-init as PID 1: signal forwarding + zombie reaping (Xvfb, x11vnc, Chrome).
ENTRYPOINT ["/usr/bin/dumb-init", "--", "/entrypoint.sh"]

EXPOSE 8080 5900
CMD ["./server"]
