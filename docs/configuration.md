# Configuration

ciso-gophish reads a JSON configuration file (`config.json` in the working directory by default; pass `--config <path>` to change it). In the container image, the environment variables below are applied to `config.json` at startup by `docker/run.sh`.

## `config.json`

```json
{
  "admin_server": {
    "listen_url": "127.0.0.1:3333",
    "use_tls": true,
    "cert_path": "gophish_admin.crt",
    "key_path": "gophish_admin.key",
    "trusted_origins": []
  },
  "phish_server": {
    "listen_url": "0.0.0.0:80",
    "use_tls": false,
    "cert_path": "example.crt",
    "key_path": "example.key"
  },
  "db_name": "sqlite3",
  "db_path": "gophish.db",
  "migrations_prefix": "db/db_",
  "contact_address": "",
  "logging": {
    "filename": "",
    "level": ""
  }
}
```

### Fields

| Field | Meaning |
|---|---|
| `admin_server.listen_url` | Address the admin console/API binds to. Keep it private (localhost, VPN, or IP-allowlisted). |
| `admin_server.use_tls` | Serve the admin console over HTTPS (recommended). |
| `admin_server.cert_path` / `key_path` | TLS material for the admin server. A self-signed pair is generated if missing. |
| `admin_server.trusted_origins` | Extra origins allowed for cross-origin admin requests (CSRF protection). Usually empty. |
| `phish_server.listen_url` | Address the phishing server binds to (what targets reach). |
| `phish_server.use_tls` | Serve landing pages over HTTPS. Provide a real certificate for public use. |
| `phish_server.cert_path` / `key_path` | TLS material for the phishing server. |
| `db_name` | `sqlite3` (default) or `mysql`. |
| `db_path` | SQLite file path, **or** the MySQL DSN when `db_name` is `mysql`. |
| `migrations_prefix` | Location prefix of the migration files (leave as default). |
| `contact_address` | Address shown to recipients who make a transparency request. |
| `logging.filename` | Log file (empty = stdout). |
| `logging.level` | Log level (empty = default). |

## Environment variables (container)

`docker/run.sh` maps these variables onto `config.json` before launching the server:

| Variable | Sets |
|---|---|
| `ADMIN_LISTEN_URL` | `admin_server.listen_url` |
| `ADMIN_USE_TLS` | `admin_server.use_tls` (JSON boolean) |
| `ADMIN_CERT_PATH` / `ADMIN_KEY_PATH` | admin TLS paths |
| `ADMIN_TRUSTED_ORIGINS` | `admin_server.trusted_origins` (comma-separated) |
| `PHISH_LISTEN_URL` | `phish_server.listen_url` |
| `PHISH_USE_TLS` | `phish_server.use_tls` (JSON boolean) |
| `PHISH_CERT_PATH` / `PHISH_KEY_PATH` | phishing TLS paths |
| `CONTACT_ADDRESS` | `contact_address` |
| `DB_NAME` | `db_name` (`sqlite3` or `mysql`) |
| `DB_FILE_PATH` | `db_path` (SQLite path or MySQL DSN) |

## Database

- **SQLite (default):** a single file at `db_path`. Persist it with a volume. Migrations run automatically on boot.
- **MySQL:** set `DB_NAME=mysql` and point `DB_FILE_PATH` at a DSN, e.g. `gophish:password@tcp(db:3306)/gophish?charset=utf8mb4&parseTime=True&loc=UTC`.

Schema changes are applied automatically at startup using the migrations under `db/db_sqlite3/migrations/` and `db/db_mysql/migrations/`.
