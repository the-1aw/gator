# Gator 🐊

A CLI for aggregating and browsing RSS feeds, backed by PostgreSQL.

## Requirements

- Go 1.25+
- A PostgreSQL Database

**These are required but provided through go tool as tool dependencies**
- [goose](https://github.com/pressly/goose) (migrations)
- [sqlc](https://sqlc.dev) (only needed if modifying SQL queries/schema)

## Setup (TBD)

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

**Running the app requires a running PostgreSQL instance.**  
You can start one with Docker if you don't have one:
 ```bash
 docker run --name gator-db \
   -e POSTGRES_USER=gator-user \
   -e POSTGRES_PASSWORD=password \
   -e POSTGRES_DB=gator \
   -p 5432:5432 \
   -v postgres-data:/bar/lib/postgresql/data \
   -d postgres
 ```

1. Create `~/.gatorconfig.json`:
   ```json
   {
     "db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable",
     "current_user_name": ""
   }
   ```

2. Run migrations:

    ```bash
        mv .env.exemple .env
        # Fill in .env with your credentials
        go tool goose up
    ```

    OR

   ```bash
   go tool goose -dir sql/schema postgres "$DB_URL" up
   ```

3. Build the binary:
   ```bash
   make build      # compile and produce ./gator binary
   ```

## Planned improvements
- [x] Refactor `cli/cli.go` into smaller packages
- [x] Add makefile command to setup and install project and update readme development
- [ ] Use XDG_CONFIG_HOME for default gator config
- [ ] Tidy migrations (since nothing runs in production we could remove altering migrations and only keep the table creation)
- [ ] Update project design so it can be database agnostic
- [ ] Update readme setup and rename it install
- [ ] Background service manager to keep `agg` running and restart on crash
- [ ] User-friendly error messages for SQL errors
- [ ] Use timestamps instead of dates for every created_at and updated_at fields
- [ ] Add default values for uuid and timestamp so we don't need to create one everytime we want to insert something from the db
- [ ] HTTP API with authentication for remote access
