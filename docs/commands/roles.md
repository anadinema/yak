# yak roles

List available role tiers.

## Usage

```sh
yak roles [flags]
```

## Description

Displays all role tiers defined in `[roles]` with their power grades, resolved role names (masked), and any configured aliases.

## Flags

| Flag               | Description                                      |
|--------------------|--------------------------------------------------|
| `--account <name>` | Show only roles available for a specific account |
| `--json`           | Output as JSON                                   |

## Example output

```
TIER       POWER GRADE  ROLE NAME (resolved)    ALIASES
read_only  100          DiagnosticsRole         ro
developer  200          DeveloperRole           dev
admin      300          PrivilegedRole    adm
root       400          GodModeRole          rt
```

## With --account prod

```
TIER       POWER GRADE  ROLE NAME (resolved)    ALIASES  OVERRIDE
read_only  100          ReaderRole            ro       yes (prod override)
developer  200          DeveloperRole           dev      —
admin      300          AdministratorRole               adm      yes (prod override)
root       400          —                       rt       not allowed on this account
```
