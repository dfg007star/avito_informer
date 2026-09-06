# Glossary

Domain terms for `avito_informer`. Keep this in sync when a term's meaning shifts.

## Link

A saved Avito search, stored in the `links` table. Has a `name`, a search-results `url`,
and an optional `min_price` / `max_price` window. The collector visits each link's `url`
and extracts the listings shown on that results page. "Parsing a link" = scraping one
results page.

## Item

A single Avito listing, stored in the `items` table, always belonging to one link
(`link_id`). Identified across scrapes by `uid` — Avito's own `data-item-id` — which is
UNIQUE, so re-scraping a link never duplicates an item.

## Cycle

One full pass of the collector over **all** links: fetch the link list, scrape each one
once (serially, in shuffled order, with a jittered delay between), then start again. The
browser context is recycled at the end of each cycle.

## Notified (`is_notify`)

Whether an item has been pushed to Telegram. Set once, on the item's first scrape, and
never re-sent. A link's *very first* scrape marks all found items as already-notified
(`link.ItemsCount == 0` at parse time) so adding a link does not flood the operator with
its entire existing inventory.

## Profile

The Chromium `--user-data-dir` directory, kept on a Docker volume. Holds the logged-in
Avito session — cookies, localStorage, IndexedDB, and PerimeterX device state. Bootstrapped
once by the operator over VNC; reused by every collector run. A "burned" profile is one
PerimeterX no longer trusts; it requires a manual VNC re-login.

## Block / Challenge

PerimeterX's "confirm you are human" interstitial. Avito serves it instead of the results
page when its bot score for the current session (IP + browser fingerprint + behaviour +
request rate) is too high. The collector treats a cycle where the results grid never
appears as blocked.

## Backoff

The collector's response to a detected block: retry the same link after 10 min, 30 min,
1 h, then stop and wait for an operator `/resume`. Distinct from `pause`, which is
operator-initiated.

## Operator

The single person running this instance. Controls the collector from Telegram (`/status`,
`/pause`, `/resume`, `/shot`, `/scrape`) and does the one-time VNC profile login.

## Recycle

Tearing down and rebuilding the `chromedp` browser **context** (one tab, its network/fetch
state) while keeping the `ExecAllocator` process and the on-disk profile. Done once per
cycle to bound Chromium's memory growth.
