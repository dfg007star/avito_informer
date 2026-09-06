# ADR 0001: Collector rework — persistent authenticated browser

Status: accepted
Date: 2026-09-06

## Context

The collector scrapes Avito search-result pages with a headless Chromium (chromedp).
Current behaviour:

- Fresh `chromedp.NewContext()` (new tab + `network`/`fetch` re-init) **per link**.
- Anonymous session — a throwaway cookie fetch at process start, reused for every scrape.
- Spoofed `navigator` properties via injected JS.

Observed result: PerimeterX ("confirm you are human") challenges roughly 18 of 20 links
per cycle. Effectively unusable. Retries do not help — the challenge is IP + fingerprint +
behaviour scored, not per-request.

Secondary problems in the current loop:

- `for {}` in `app.collect` never checks `ctx.Done()` and only `time.Sleep`s *inside* the
  per-link loop, so an empty `links` table produces a busy-spin at ~80 % CPU.
- Errors from `GetAllLinks` / `CreateItems` propagate to `main`, where the return value of
  `fmt.Errorf(...)` is discarded and the process exits 0. Docker restarts it silently.
- `parser.Parse` indexes `imageUrls[0]` unconditionally — panics on any item with no images.
- `parser.Shutdown()` is never called.

## Decision

Rework the collector into a **long-running service driven by a single persistent,
logged-in browser**, scraping **serially** at a deliberately slow, jittered cadence, with
**Telegram** as the operator control channel.

### Browser

- One `chromedp` `ExecAllocator` **and one browser context** for the whole process.
- Chromium runs **non-headless under Xvfb** inside the container.
- `--user-data-dir` points at a Docker volume (`avito-profile`). The operator logs into
  Avito **once**, by hand, over VNC into the container. The profile (cookies, localStorage,
  IndexedDB, PerimeterX device state) then survives restarts.
- Each link = `chromedp.Navigate` in the **same tab**. No new context per link.
- The context is **recycled once per full cycle** (after every link has been visited once):
  tear down + rebuild the context, keep the allocator and the on-disk profile. Bounds the
  Chromium memory leak at the cost of one ~3 s re-spawn per cycle.
- Per-navigation timeout retained (wedged navigate cancels without killing the process).

### Cadence

- **Serial.** One link at a time. Loop over all links, then start a new loop.
- Target ~30 s per link. Actual delay **randomised per link** (e.g. 20–45 s) and **link
  order shuffled** each cycle, to avoid a fixed scraping signature.
- No proxy. One Avito account, one instance. Replicas with separate accounts + proxies are
  a later, out-of-scope step.

### Block handling

- Detect a PerimeterX challenge / non-loading results page. Signal: the results grid
  (`div[data-marker='item']`) does not appear within the navigation timeout **and/or** the
  page shows a known block marker (exact markers TBD from a captured block page — tracked as
  a `ponytail:` in code until confirmed).
- On block: mark state `blocked`, send a Telegram alert, and **back off** — retry the same
  link after 10 min, then 30 min, then 1 h, then stop and wait for `/resume`.
- While `paused` or stopped: do nothing until an operator command.

### Telegram control

- `notification` **owns the bot** (it already runs one). Telegram is the single control
  surface — no public REST API on the collector.
- Operator commands (served by `notification`, which calls a small **internal-only** HTTP
  endpoint on `collector`; not a public API):

  | Command      | Effect                                                        |
  |--------------|--------------------------------------------------------------|
  | `/status`    | last cycle time, links parsed N/M, blocked?, logged-in?, RAM |
  | `/pause`     | stop scraping after the current link                         |
  | `/resume`    | resume; also clears a backoff and retries the last link      |
  | `/shot`      | screenshot of the current browser tab                        |
  | `/scrape <id>` | force-scrape one link now                                  |

- The collector also exposes a **read-only dashboard** on its existing `:8080` (currently
  `EXPOSE`d, unused) and via the `http` service UI: current status, last cycle, recent logs,
  and a form to paste fresh cookies / trigger a profile reload.

### Notifications

- `is_notify` semantics unchanged: a listing is notified **once**, ever.
- Price-drop / re-notify on `min_price`/`max_price` window changes is **future work**, not
  in this rework.

## Phasing

1. **Loop hardening** — `ctx.Done()`, always-sleep, propagate errors instead of silent
   exit, guard `imageUrls[0]`, call `Shutdown`. Ships on its own.
2. **Persistent browser + profile** — one context, `UserDataDir`, Xvfb + VNC in the image,
   reuse-tab scraping, per-cycle recycle. Bootstrap doc.
3. **Block detection + backoff + Telegram alert.**
4. **Telegram command channel** + collector internal endpoint + `http` dashboard.
5. *(later, out of scope)* proxy per instance, multiple accounts.

## Consequences

- Collector image grows (Xvfb, x11vnc, a window manager). RAM per instance rises; the
  `mem_limit: 2g` / `shm_size: 2g` in compose stay and are the ceiling.
- `notification` gains a runtime dependency on `collector` being reachable for the command
  path. The scrape loop itself stays DB-only and works if `notification` is down.
- The approach is coupled to PerimeterX's current behaviour. If Avito changes the
  challenge, Phase 2/3 may need revisiting; Phases 1 and the service structure survive.
- A burned profile requires a manual VNC re-login. Expected to be rare (operator reports
  staying logged in for a week+ from a desktop); if it becomes frequent the approach fails
  and we revisit (official API, solver service).

## Open items

- Capture a real PerimeterX block page (HTML + screenshot) to pin the block-detection
  markers. Until then detection is "results grid absent after timeout".
- Confirm desktop-Chrome vs container-Chromium profile compatibility is **not** relied on
  (decision: bootstrap over VNC *inside* the container, so versions always match).
