# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

---

## Installation

```bash
go install github.com/youruser/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and log file:

```bash
portwatch start --interval 60 --log /var/log/portwatch.log
```

Run a one-time snapshot of currently open ports:

```bash
portwatch scan
```

**Example output:**

```
[INFO]  Baseline established: 3 open ports (22, 80, 443)
[ALERT] New port detected: 8080 (tcp)
[INFO]  Port closed: 80 (tcp)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `30` | Scan interval in seconds |
| `--log` | stdout | Path to log output file |
| `--config` | `~/.portwatch.yaml` | Path to config file |

---

## License

MIT © 2024 youruser