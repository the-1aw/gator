# Gator 🐊

A CLI for aggregating and browsing RSS feeds, backed by PostgreSQL.

## Requirements

- Go 1.25+
- PostgreSQL
- [goose](https://github.com/pressly/goose) (migrations)
- [sqlc](https://sqlc.dev) (only needed if modifying SQL queries/schema)

## Setup

1. Create a PostgreSQL database.

2. Create `~/.gatorconfig.json`:
   ```json
   {
     "db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable",
     "current_user_name": ""
   }
   ```

3. Run migrations:
   ```bash
   goose -dir sql/schema postgres "$DB_URL" up
   ```

4. Build and install:
   ```bash
   go install .
   ```

## Usage

```bash
gator <command> [arguments]
```

### Commands

| Command | Arguments | Description |
|---|---|---|
| `register` | `<name>` | Create a new user and switch to them |
| `login` | `<name>` | Switch to an existing user |
| `users` | | List all users |
| `reset` | | Delete all users and data |
| `addfeed` | `<name> <url>` | Add an RSS feed and follow it |
| `feeds` | | List all feeds |
| `follow` | `<url>` | Follow an existing feed |
| `following` | | List feeds you follow |
| `unfollow` | `<url>` | Unfollow a feed |
| `agg` | `<interval>` | Start polling feeds (e.g. `30s`, `5m`) |
| `browse` | `[limit]` | Print recent posts (default: 2) |

### Example workflow

```bash
gator register alice
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
gator agg 1m        # runs in foreground, fetching feeds every minute
gator browse 10     # in another terminal, browse the latest posts
```

## Development

```bash
go build ./...
go test ./...

# After editing sql/queries/ or sql/schema/
sqlc generate

# Run a migration rollback
goose -dir sql/schema postgres "$DB_URL" down
```

## Planned improvements
- [x] Refactor `cli/cli.go` into smaller packages
- [ ] Use timestamps instead of dates for every created_at and updated_at fields
- [ ] Add default values for uuid and timestamp so we don't need to create one everytime we want to insert something from the db
- [ ] User-friendly error messages for SQL errors
- [ ] HTTP API with authentication for remote access
- [ ] Background service manager to keep `agg` running and restart on crash
