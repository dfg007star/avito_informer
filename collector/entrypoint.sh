#!/usr/bin/env bash
# Starts Xvfb (always) and x11vnc + fluxbox (when VNC_ENABLED=true), then execs
# the collector. Chromium needs a real display for non-headless operation;
# PerimeterX flags headless Chromium.
set -e

: "${DISPLAY:=:99}"
export DISPLAY

# Clear a stale X lock left by an unclean shutdown, else Xvfb refuses to start.
rm -f "/tmp/.X${DISPLAY#:}-lock" "/tmp/.X11-unix/X${DISPLAY#:}" 2>/dev/null || true

# Clear stale Chromium profile singleton locks (unclean shutdown leaves them and
# every subsequent launch fails "profile appears to be in use").
: "${PROFILE_DIR:=/data/avito-profile}"
rm -f "${PROFILE_DIR}/SingletonLock" "${PROFILE_DIR}/SingletonCookie" \
      "${PROFILE_DIR}/SingletonSocket" 2>/dev/null || true

# Chromium on Alpine is unstable without a D-Bus session bus.
mkdir -p /run/dbus
dbus-daemon --system --fork 2>/dev/null || true
eval "$(dbus-launch --sh-syntax)" 2>/dev/null || true

Xvfb "$DISPLAY" -screen 0 1366x768x24 -nolisten tcp &
XVFB_PID=$!

# wait for the display socket
for _ in $(seq 1 50); do
  [ -e "/tmp/.X11-unix/X${DISPLAY#:}" ] && break
  sleep 0.1
done

if [ "${VNC_ENABLED}" = "true" ]; then
  fluxbox >/dev/null 2>&1 &
  if [ -n "${VNC_PASSWORD}" ]; then
    x11vnc -storepasswd "${VNC_PASSWORD}" /tmp/.vncpass >/dev/null 2>&1
    x11vnc -display "$DISPLAY" -forever -shared -rfbauth /tmp/.vncpass -rfbport 5900 -bg >/dev/null 2>&1
  else
    echo "WARNING: VNC_ENABLED=true with no VNC_PASSWORD — desktop is open without auth"
    x11vnc -display "$DISPLAY" -forever -shared -nopw -rfbport 5900 -bg >/dev/null 2>&1
  fi
  echo "VNC listening on :5900"
fi

term() {
  kill "$XVFB_PID" 2>/dev/null || true
  exit 0
}
trap term SIGTERM SIGINT

exec "$@"
