# Stockyard Steward

**Self-hosted expense tracker.** Single Go binary, embedded SQLite, no external dependencies.
Part of the [Stockyard](https://stockyard.dev) suite of self-hosted developer tools.

Submit expenses, attach receipts, approve or reject, export to CSV.

## Install

```bash
curl -fsSL https://stockyard.dev/steward/install.sh | sh
```

That downloads the latest release for your platform from
[GitHub releases](https://github.com/stockyard-dev/stockyard-steward/releases/latest)
and drops a single binary on disk. Linux (amd64/arm64), macOS (Intel/Apple
Silicon), and Windows (amd64) are all supported.

Then run it:

```bash
stockyard-steward
```

Dashboard at [http://localhost:9840/ui](http://localhost:9840/ui).
HTTP API at `http://localhost:9840/api`. Data lives in `~/.stockyard/steward/`
by default — override with `-data <path>` or `DATA_DIR=<path>`.

## Personalization

Steward is one of 169 tools in the Stockyard suite. When you install it via
the [Stockyard launcher](https://github.com/stockyard-dev/stockyard-launcher) or
through the AI toolkit builder at [stockyard.dev](https://stockyard.dev), the
launcher writes a `config.json` into the data directory that personalizes the
dashboard for your specific use case — custom field labels, default values,
terminology that matches your business.

Without `config.json`, the tool runs with sensible defaults for the
"self-hosted expense tracker" category. With it, the same binary serves a dashboard
tailored to your business in 30 seconds, no rebuild required. Read the
config schema at `GET /api/config`.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9840` | HTTP listen port |
| `DATA_DIR` | `~/.stockyard/steward` | SQLite + config directory |
| `STOCKYARD_LICENSE_KEY` | *(empty)* | Pro license key (see below) |

Command-line flags `-port` and `-data` override the env vars.

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Use case | Personal, hobby, small teams | Production, paid customers |
| Limits | Per-tool free tier (see [pricing](https://stockyard.dev/pricing/)) | Unlimited |
| Price | $0/mo | $0.99/mo per tool, or $7.99/mo for the full Stockyard suite |

Get a Pro license at [stockyard.dev/pricing/?plan=steward-pro](https://stockyard.dev/pricing/?plan=steward-pro).
14-day free trial, cancel anytime, your data stays on your machine either way.

## Screenshot

See the live dashboard and feature tour at
[stockyard.dev/steward/](https://stockyard.dev/steward/).

## License

BSL 1.1 — free for non-production use, free for production use under the
Stockyard Pro license, converts to Apache 2.0 four years after release.
See `LICENSE` for the full terms.
