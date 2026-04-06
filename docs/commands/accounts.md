# yak accounts

List all configured accounts.

## Usage

```sh
yak accounts [flags]
```

## Description

Displays all accounts defined in your config with their IDs (partially masked), default role, allowed roles, and whether safeguards are active.

## Flags

| Flag     | Description    |
|----------|----------------|
| `--json` | Output as JSON |

## Example output

```
NAME      ID               DEFAULT ROLE  ALLOWED ROLES              SAFEGUARD
dev       123456****012    admin         all                        —
sandbox   234567****123    admin         all                        —
test      345678****234    admin         all                        —
prod      456789****345    developer     read_only, developer, admin active
spt       567890****456    developer     read_only, developer, admin active
spp       678901****567    developer     read_only, developer, admin active
mgmt      789012****678    admin         read_only, developer, admin —
```
