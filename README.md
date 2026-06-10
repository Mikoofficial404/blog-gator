# Gator 🐊

Gator is a multi-user RSS feed aggregator CLI built with Go and PostgreSQL. It lets you follow RSS feeds, automatically fetch posts in the background, and browse the latest content right in your terminal.

## Prerequisites

Before running Gator, make sure you have the following installed:

- [Go](https://golang.org/dl/) (1.21 or later)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

Install the CLI using `go install`:

```bash
go install github.com/Mikoofficial404/blog-go@latest
```

Make sure your Go binary path is in your `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Configuration

Create a config file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator"
}
```

Replace `username`, `password`, and `gator` with your PostgreSQL credentials and database name.

### Database Setup

Run the migrations to set up the database schema:

```bash
goose -dir ./sql/schema postgres "your_connection_string" up
```

## Usage

### Register a new user

```bash
gator register <username>
```

### Login as an existing user

```bash
gator login <username>
```

### Add a feed

```bash
gator addfeed "Feed Name" "https://feed-url.com/rss"
```

Example:

```bash
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
```

### List all feeds

```bash
gator feeds
```

### Follow a feed

```bash
gator follow "https://feed-url.com/rss"
```

### Unfollow a feed

```bash
gator unfollow "https://feed-url.com/rss"
```

### List followed feeds

```bash
gator following
```

### Start the aggregator

Run this in a separate terminal to continuously fetch posts in the background:

```bash
gator agg 1m
```

The argument is the time between requests (e.g. `30s`, `1m`, `1h`).

### Browse posts

```bash
gator browse        # shows 2 posts (default)
gator browse 10     # shows 10 posts
```

### List all users

```bash
gator users
```

### Reset the database

```bash
gator reset
```

## Typical Workflow

1. Register and login
2. Add some RSS feeds with `addfeed`
3. Run `gator agg 1m` in one terminal to start fetching posts
4. Use `gator browse` in another terminal to read the latest posts
