# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`avito_informer` scrapes Avito (avito.ru) search-result pages for saved searches ("links"),
stores new listings in Postgres, and pushes new listings to Telegram. Three Go services plus a
shared platform module, wired together only through the Postgres database — there is no RPC
between services.

## Repository layout

Go multi-module repo. Each service is its own module with its own `go.mod`; they are joined at
build time by a `go.work` file that is **git-ignored and generated** (see `deploy/compose/core/run.sh`
and both Dockerfiles). To work locally you must create it:

```sh
go work init
go work use ./collector ./http ./notification ./platform
go work sync
```

Modules:

- `http/` — web UI (server-rendered Go `html/template`) to manage saved searches and view scraped items. Owns the DB migrations.
- `collector/` — headless-Chrome scraper (chromedp). Long-running loop, no HTTP server.
- `notification/` — polls the DB for un-notified items and sends them to Telegram. Also runs a Telegram bot for `/start`.
- `platform/` — shared library only. Currently just `pkg/migrator` (goose wrapper).

Module import path base: `github.com/dfg007star/avito_informer/<module>`.

## Git

- Do **not** add `Co-Authored-By` or `Claude-Session` trailers to commit messages, and do
  not add "Generated with Claude Code" to PR descriptions. Plain messages only.
- Commit author identity is set in local `.git/config` only (`dfg007star` /
  `lastfm4321@gmail.com`) — never touch global git config.
- `origin` is `github.com/dfg007star/avito_informer`. Pushes authenticate with a fine-grained
  PAT supplied per-session; it is never written to disk or `.git/config`.

## Commands

Task runner is [Task](https://taskfile.dev). All `task` targets operate on Docker Compose, not local Go:

- `task run` / `task stop` — start/stop the compose stack
- `task rebuild` — `build --no-cache` + `up -d` (no git pull)
- `task update` — `git reset --hard origin/main` then rebuild (deploy-box target; destructive to local changes)
- `task logs` — tail all services
- `task deps:update` — `go mod tidy` across `http collector platform` (note: `notification` is missing from the `MODULES` var)

Compose file: `deploy/compose/core/docker-compose.yml`. Needs `deploy/compose/core/.env`
(copy from `.env.example`). `deploy/compose/core/load_system_deps.sh` provisions a bare Ubuntu box
(Docker, Go, Task, Fish).

Build/test/lint a single module (run inside the module dir, needs `go.work`):

```sh
cd collector
go build ./...
go test ./...
go test ./internal/parser/ -run TestName -v
```

`golangci-lint`, `gci`, `gofumpt` binary paths are declared in `Taskfile.yml` vars but no task
invokes them and there is no config file — treat lint as gofumpt + gci + golangci-lint if you add it.

There are currently **no test files** in the repo.

### Running a service outside Docker

`main.go` in each service checks `RUNNING_IN_DOCKER`. When not `"true"` it loads config from
`../deploy/compose/core/.env.local` — so run from inside the module directory (`cd http && go run ./cmd/main.go`).

## Architecture

### Data model (owned by `http/migrations/`, goose format)

- `links` — a saved Avito search: `name`, `url`, optional `min_price`/`max_price`, `parsed_at`.
- `items` — a scraped listing: `link_id` FK (ON DELETE CASCADE), `uid` (Avito's `data-item-id`, UNIQUE),
  `title`, `price`, `url`, `preview_url`, `is_notify` (has it been pushed to Telegram yet).

Migrations run automatically on `http` service startup (`app.initMigrator`), using a second
`database/sql` connection opened from the pgx config. The other services do **not** migrate.

### Service flow

```
http  (user adds a link) ─────▶ links table
collector  loop: read all links ─▶ chromedp scrape each ─▶ insert new items (dedup by uid)
notification  loop: SELECT items WHERE is_notify=false ─▶ send Telegram ─▶ UPDATE is_notify=true
```

- **collector** (`internal/app/app.go` `collect`): infinite `for` loop, no ctx cancellation check.
  Gets a set of Avito cookies once at startup via a throwaway URL, reuses them for every scrape.
  `internal/parser/parser.go` is the scraper: one `chromedp.ExecAllocator` for the process, a fresh
  browser context per link, image requests blocked via `fetch` domain, anti-bot JS injected via
  `AddScriptToEvaluateOnNewDocument`, retries `GET_ITEMS_MAX_RETRY` times, parses the rendered HTML
  with goquery selectors keyed off Avito's `data-marker` / `itemprop` attributes (these break when
  Avito changes markup — the CSS class `photo-slider-list-R0jle` is especially fragile).
  New item gets `is_notify = (link had 0 items before)` so the first scrape of a new link doesn't
  spam every existing listing.
- **notification** (`internal/app/app.go` `Run`): 10s poll loop, batch of 30 newest un-notified,
  renders `internal/service/telegram/templates/item_notification.tmpl` (embedded), sends HTML-parse-mode
  message, marks `is_notify=true` one row at a time. Telegram bot runs in a background goroutine.
- **http**: `net/http.ServeMux` with Go 1.22 method patterns. All routes except `/login` and `/static/`
  go through `handler/auth.Middleware`, which is a single shared password (`HTTP_PASSWORD`) stored
  verbatim in an `auth` cookie. The UI is htmx-style partial templates (`html/_links_*.html`).

### Per-service internal package convention

Every service repeats the same layered structure:

```
internal/
  app/        App struct + hand-rolled DI container (di.go), lazy singletons, panic on wiring failure
  config/     interfaces.go (consumer interfaces) + env/ (caarlos0/env structs, one per concern)
  model/      domain structs
  service/    service.go (interface) + <name>/ (impl split one-file-per-method)
  repository/ repository.go (interface) + <name>/ (impl) + model/ (DB row structs) + converter/ (row↔domain)
  handler/ or client/  transport layer
```

- **Config**: `config.Load(path...)` is called once from `main`, stores a package-global
  `appConfig`, read everywhere via `config.AppConfig()`. Env vars are `,required` — a missing var
  is a hard startup failure. Each `env/*.go` file is a private struct with getter methods; the
  matching interface lives in `config/interfaces.go`.
- **DI**: `diContainer` in `internal/app/di.go`. No framework — `Xxx(ctx)` methods with
  `if d.x == nil` memoization. Postgres client is a single non-pooled `pgx.Conn` per process.
- **DB access**: `Masterminds/squirrel` query builder with `PlaceholderFormat(squirrel.Dollar)`,
  executed through `pgx`. `gorm` is a dependency but only its types are used (not the ORM).
  Repos return domain models; `converter/` translates between `repository/model` row structs and
  `internal/model` domain structs.

## Conventions & gotchas

- Errors are frequently created with `fmt.Errorf(...)` and then **not returned or logged** —
  the value is discarded (see every `main.go`, `app.go` `runHTTPServer`, etc.). If you touch that
  code, actually return/handle the error rather than mirroring the existing bug.
- `notification` module is omitted from `Taskfile.yml`'s `MODULES`/`SERVICES` vars and from
  `platform`'s replace wiring — check when adding cross-module deps.
- `platform` is published as a pseudo-version in the other modules' `go.mod` but resolved locally
  via `go.work` for dev and via `COPY` in the Dockerfiles.
- Docker builds copy **all four** module dirs into the image (needed for `go work sync`), then
  build one service.
- Collector compose has `collector` and `collector-proxy-1` (same image, `PROXY` env) — multiple
  scraper instances are expected; the parser reads `PROXY` but does not currently apply it to chromedp.
- Secrets: `.env`, `.env.local` are git-ignored. `HTTP_PASSWORD` is the entire auth system.
