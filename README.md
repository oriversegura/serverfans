# ServerFans

A terminal UI for controlling server fan speeds via IPMI, built with Go and [Charmbracelet Huh](https://github.com/charmbracelet/huh).

## Requirements

| Requirement | Notes |
|---|---|
| `ipmitool` | Must be installed and available in your PATH |
| Network access | Same network segment as the target server |
| IPMI credentials | User with **Operator** or **Administrator** privileges |

## Installation

```bash
git clone https://github.com/oriversegura/serverfans
cd serverfans
go build -o serverfans .
```

## Usage

**Interactive mode** — prompts for all values:

```bash
./serverfans
```

**Config file mode** — pre-fills fields from a JSON file:

```bash
./serverfans config.json
```

The config file can provide any combination of fields. Missing fields are prompted interactively.

### Config file format

```json
{
  "ip":       "192.168.1.50",
  "user":     "admin",
  "password": "secret",
  "speed":    "40"
}
```

> **Tip:** Omit `password` from the file and enter it interactively to avoid storing credentials on disk.

## Fan speed range

| Value | Meaning |
|---|---|
| `10` | Minimum (10%) |
| `100` | Maximum (100%) |

## How it works

ServerFans sends two raw IPMI commands to the target host:

1. **Enable manual fan control** — disables the BMC's automatic fan algorithm.
2. **Set fan speed** — applies the requested percentage to all fans.

```
ipmitool -I lanplus -H <ip> -U <user> -P <pass> raw 0x30 0x30 0x01 0x00
ipmitool -I lanplus -H <ip> -U <user> -P <pass> raw 0x30 0x30 0x02 0xff <hex_speed>
```

> These raw commands target Dell iDRAC. Other vendors (HP iLO, Supermicro) may use different OEM commands.

## Project structure

```
serverfans/
├── serverfans.go          # Entry point — orchestrates the run loop
├── internal/
│   ├── config/
│   │   └── config.go      # Config struct and JSON loader
│   ├── ipmi/
│   │   └── client.go      # ipmitool wrappers (CheckInstalled, SetFanSpeed)
│   └── ui/
│       └── ui.go          # Banner, forms, and styled output
├── example_config.json
└── go.mod
```

## License

MIT
