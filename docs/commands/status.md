# yak status

Show current active account, role tier, and login state.

## Usage

```sh
yak status [flags]
```

## Description

Reads the state file and displays a summary of your current yak context.

## Flags

| Flag              | Description                                                        |
|-------------------|--------------------------------------------------------------------|
| `--json`          | Output as JSON instead of formatted text                           |
| `--print-profile` | Print only the active AWS profile name (used by the shell wrapper) |

## Example output

```
Account:     prod
Role:        developer
Region:      eu-west-1
Last login:  2026-04-06 10:30:00 UTC (2 hours ago)
Safeguard:   active (max: developer)
```

## JSON output

```json
{
  "account": "prod",
  "role": "developer",
  "region": "eu-west-1",
  "last_login": "2026-04-06T10:30:00Z",
  "safeguard_active": true,
  "safeguard_max_role": "developer"
}
```
