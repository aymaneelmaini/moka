# 🏦 Moka - Personal Finance Tracker

- personal project to use by myself to track my money, and explore **go-lang** + **HTMX** stack at the same time.
- used database is **Sqlite** because the project is simple no need for extra setup and work
- Project name: **Moka** means money in moroccan dialect

A clean, modern personal finance tracker built with **Go + HTMX + SQLite** following **Clean Architecture**, **Domain-Driven Design (DDD)**, and **Functional Programming** principles.

## Installation

one-line install:
```bash
curl -fsSL https://raw.githubusercontent.com/aymaneelmaini/moka/main/install.sh | bash
```

this downloads the binary, installs as systemd service, sets up local domain (moka.local), and runs in background.

access at: http://moka.local:9876

manage service:
```bash
sudo systemctl status moka    # check status
sudo journalctl -u moka -f    # logs
sudo systemctl restart moka   # restart
```

data stored in `~/.moka/moka.db`

