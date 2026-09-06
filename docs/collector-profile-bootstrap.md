# Bootstrapping the Avito profile (one-time)

The collector scrapes with a **persistent, logged-in Chromium profile** stored in the
`avito_profile` Docker volume (mounted at `/data/avito-profile`). You seed it once by
logging into Avito by hand over VNC. After that every collector restart reuses the session.

Repeat this only when the profile is "burned" (collector reports `not logged in` /
persistent blocks that a `/resume` doesn't clear).

## Steps

1. In `deploy/compose/core/.env` set:

   ```
   COLLECTOR_VNC_ENABLED=true
   VNC_PASSWORD=<pick-something>
   PARSER_HEADLESS=false
   ```

2. Recreate the collector:

   ```
   docker compose -f deploy/compose/core/docker-compose.yml up -d collector
   ```

   The scrape loop still runs, but you'll see `WARNING: Avito session is not logged in`
   until step 5. That's expected.

3. From your machine, tunnel to the VNC port (it's bound to `127.0.0.1` on the Docker host):

   ```
   ssh -L 5900:127.0.0.1:5900 <docker-host>
   ```

   Skip this if you're already on the host.

4. Connect a VNC client to `localhost:5900`, password = `VNC_PASSWORD`.
   You get a bare fluxbox desktop. The collector's Chromium window is there (right-click
   the desktop → clients, or it's already focused).

5. In that Chromium: go to `https://www.avito.ru`, log in with the Avito account, solve
   the PerimeterX puzzle if shown. Confirm you can open a search-results page and see
   listings. Leave the browser as-is — do not close it.

6. Turn VNC back off:

   ```
   COLLECTOR_VNC_ENABLED=false
   ```

   ```
   docker compose -f deploy/compose/core/docker-compose.yml up -d collector
   ```

   On restart the collector should log `Avito session is logged in`.

## Notes

- The profile lives in the `avito_profile` volume. `docker compose down` keeps it;
  `docker compose down -v` **deletes** it and you start over.
- Only one Chromium can use the profile dir at a time — the collector holds it. Don't run
  a second browser against the same volume.
- VNC port is `127.0.0.1`-only by design. Never expose 5900 publicly; `x11vnc` auth is weak.
- If Chromium isn't visible in VNC, it may have crashed — check `docker logs core-collector-1`
  for `failed to start browser`.
