# mnémo
A local-first, privacy-friendly CLI to capture and search your shell history on Ubuntu/WSL (bash). Built in Go with SQLite + FTS5 for fast search
## Install

- Go **1.22+**
- SQLite is embedded via `modernc.org/sqlite` (no CGO)

```bash
go install github.com/yourname/mnemo/cmd/mnemo@latest
mnemo init
```

> WSL: keep the DB under Linux `$HOME` (e.g., `/home/<user>/.local/share/mnemo`).

## Bash hook (WSL)

Add to your `~/.bashrc`:

```bash
mnemo_preexec() {
  local ts=$(date +%s%3N)
  local cmd="$BASH_COMMAND"
  mnemo record --cmd "$cmd" --cwd "$PWD" --shell bash --ts "$ts" >/dev/null 2>&1 &
}
PROMPT_COMMAND='mnemo_preexec'
```

Then reload:

```bash
source ~/.bashrc
```

## Quick start

Type a few commands in your terminal, then:

```bash
mnemo search "docker build" --limit 10
mnemo stats
```

If you omit the query, `search` returns the most recent commands:

```bash
mnemo search --limit 20
```

## Config

Defaults → `~/.config/mnemo/config.yaml` → flags.

```yaml
db_path: ~/.local/share/mnemo/history.db
default_limit: 20
busy_timeout_ms: 2000
redact:
  enabled: true
  replacement: "[REDACTED]"
  key_names: ["password","token","api_key","secret"]
  headers: ["authorization","x-api-key"]
  url_credentials: true
```

## Positioning

- **Local-first**: all data stays on your machine  
- **Privacy-friendly**: secrets redacted before storage  
- **Fast**: optimized for ≤ 1 M commands  
- **Reliable**: WAL + busy timeout

## Limitations (MVP)

- No full fuzzy search (prefix & filters only)  
- No multi-machine sync  
- Secrets redacted via rules (may need tuning)

## License

Apache License 2.0 — see [LICENSE](./LICENSE).
