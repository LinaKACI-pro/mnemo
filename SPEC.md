# MVP Specification (lean) — mnemo

**Goal:** A **local-first**, **privacy-friendly** CLI to capture and search shell history on **Ubuntu (WSL)** with **bash**.  
**Stack:** Go (**Cobra + koanf**), **SQLite + FTS5** with **WAL**, **naive secret redaction**, **no daemon**, **no spool** (for MVP).

## 1) Scope

- Install & use directly (**one binary + one bash hook**).
- **Capture** each command (short-lived background process).
- **Fast search** (FTS), **minimal stats**.
- **Security by design**: redact before storage, no telemetry.
- Target **bash/WSL** first.

**Out of scope (MVP):** daemon, spool/flush, TUI, sync, full fuzzy, sessions/flows, doctor.

## 2) Architecture (no daemon)

```
Terminal(bash) -> Hook(PROMPT_COMMAND) -> mnemo record (bg) -> SQLite(WAL)
                                      ^                         ^
                           mnemo search / stats  --------------
```

## 3) Components

- **CLI (Cobra):** `init`, `ingest`, `record`, `search`, `stats`.
- **Config (koanf):** **defaults → config file → flags** (no ENV in MVP).
- **DB:** SQLite + FTS5, WAL + `busy_timeout`.
- **Redaction:** 3–5 simple rules, applied **before** insert.

## 4) Data model (SQL)

```sql
CREATE TABLE commands(
  id      INTEGER PRIMARY KEY,
  ts      INTEGER NOT NULL,  -- epoch ms
  shell   TEXT,              -- "bash"
  cmdline TEXT NOT NULL,     -- redacted
  cwd     TEXT
);

CREATE VIRTUAL TABLE commands_fts USING fts5(
  cmdline, content='commands', content_rowid='id',
  tokenize='unicode61 remove_diacritics 2', prefix='2 3'
);

CREATE TRIGGER commands_ai AFTER INSERT ON commands BEGIN
  INSERT INTO commands_fts(rowid, cmdline) VALUES (new.id, new.cmdline);
END;

CREATE INDEX idx_commands_ts ON commands(ts DESC);
```

**PRAGMA at `init`:**
```sql
PRAGMA journal_mode=WAL;
PRAGMA busy_timeout=2000;
```

## 5) CLI — commands & examples

- `mnemo init`  
  Creates app dir, DB, PRAGMAs, migrations.

- `mnemo ingest --bash ~/.bash_history`  
  Imports bash history (use “now” as `ts` if unknown).

- `mnemo record --cmd "<cmd>" --cwd "$PWD" --shell bash --ts 1724...`  
  Inserts **after redaction**.  
  **Retry:** 3–5 attempts (10/20/40/80/160 ms + jitter). If still locked → **drop** silently (MVP).

- `mnemo search "<query>" [--limit N] [--cwd PATH] [--since DURATION]`  
  FTS top-N; if `query` omitted, return **N most recent**.

- `mnemo stats`  
  Prints `Total: <N>` (global count).

## 6) Bash hook (WSL)

```bash
# ~/.bashrc
mnemo_preexec() {
  local ts=$(date +%s%3N)
  local cmd="$BASH_COMMAND"
  mnemo record --cmd "$cmd" --cwd "$PWD" --shell bash --ts "$ts" >/dev/null 2>&1 &
}
PROMPT_COMMAND='mnemo_preexec'
```

> **WSL tip:** keep the DB under Linux **$HOME** (e.g., `/home/<user>/.local/share/mnemo`) — avoid `/mnt/c`.

## 7) Config (koanf)

**Resolution:** defaults → `~/.config/mnemo/config.yaml` → flags.

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
max_index_bytes: 4096   # cut scanning/indexing beyond this (perf/safety)
```

## 8) Redaction (naive)

- **v0 rules** (case-insensitive):
  - `key=value` for `password|token|api_key|secret` → mask value.
  - Flags: `--password` / `--token` (both `--x=V` and `--x V`).
  - HTTP headers: `Authorization: Bearer …`, `X-Api-Key: …` → mask value.
  - URL credentials: `user:pass@host` → mask `user:pass@`.
- **Replacement:** `[REDACTED]`.
- **Cheap pre-filter:** only run regex if the line includes a keyword (`password`, `token`, `Authorization`, `://`).

## 9) Target performance (ballpark)

- `record` (DB OK): **~1–5 ms** (spawn + redact + INSERT WAL).
- `search` top-10 (≤100k rows): **≤50 ms** (often **<10 ms**).
- Ingest (bash): internal batching → **tens of k/s**.

## 10) Security

- **Local-first:** no network I/O.
- **Permissions:** app dir **0700**, DB **0600**.
- **Option:** `--no-store` to skip a sensitive command.

## 11) Included "nice-to-have"

- **README (minimal):** install, hook, 2 examples, Positioning/Limitations.
- **LICENSE:** Apache-2.0.
- **Tests:**
  - **Redaction (golden):** ~6 cases (Authorization, key=value, URL creds, flags…).
  - **Smoke insert/search:** init → insert dummy → search top-1.

## 12) Positioning & Limitations (for README)

```
Positioning
- Local-first • Privacy-friendly • Fast (≤ 1 M commands) • Reliable (WAL)

Limitations (MVP)
- No full fuzzy search (prefix & filters only)
- No multi-machine sync
- Secrets redacted via rules (may need tuning)
```

**Done criteria (v0.1):**  
Build via `go install …@latest`; `init/ingest/record/search/stats` work; bash hook works on WSL; redaction v0 on by default; README + LICENSE + 2 tests pass.
